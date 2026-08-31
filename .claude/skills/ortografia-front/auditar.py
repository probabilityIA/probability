#!/usr/bin/env python3
"""Auditoria de cobertura: busca tildes faltantes que el diccionario del skill no
cubre, contrastando contra un diccionario hunspell real de espanol.

Requiere:  pip install spylls
           curl -sSLO https://raw.githubusercontent.com/wooorm/dictionaries/main/dictionaries/es/index.dic
           curl -sSLO https://raw.githubusercontent.com/wooorm/dictionaries/main/dictionaries/es/index.aff
(los dos index.* van en esta misma carpeta, estan en .gitignore)

No corrige nada: imprime candidatos para revisar a mano y, si son reales,
agregar a diccionario.tsv. Buena parte de lo que reporta son falsos positivos
(palabras tecnicas o en ingles que el diccionario espanol "arregla": min, max,
items, COD, ROI). Se revisan uno por uno.
"""
import itertools, os, re, sys, unicodedata
sys.path.insert(0, ".claude/skills/ortografia-front")
import revisar
from spylls.hunspell import Dictionary

SP = os.path.dirname(os.path.abspath(__file__))
d = Dictionary.from_files(os.path.join(SP, "index"))
WORD = re.compile(r"(?<![A-Za-z0-9_\\])[A-Za-zÁÉÍÓÚÜÑáéíóúüñ]{3,}(?![A-Za-z0-9_])")
UESC = re.compile(r"\\u00([0-9a-fA-F]{2})")
VOW = "aeiouAEIOU"
ACC = {"a": "á", "e": "é", "i": "í", "o": "ó", "u": "ú",
       "A": "Á", "E": "É", "I": "Í", "O": "Ó", "U": "Ú"}

def accent_variants(w):
    pos = [i for i, c in enumerate(w) if c in VOW]
    for i in pos:
        yield w[:i] + ACC[w[i]] + w[i + 1:]

paths = ["front/central/src", "front/website/src"]
revisar.collect_protected(paths)
hits = {}
for path in revisar.walk(paths):
    try:
        lines = open(path, encoding="utf-8").read().split("\n")
    except (UnicodeDecodeError, OSError):
        continue
    for n, line in enumerate(lines, 1):
        if revisar.LINE_SKIP.search(line):
            continue
        for kind, off, text in revisar.spans(line):
            if len(text.split()) < 2:
                continue
            text = UESC.sub(lambda mm: chr(int(mm.group(1), 16)), text)
            for m in WORD.finditer(text):
                w = m.group(0)
                if any(ord(c) > 127 for c in w):
                    continue
                if d.lookup(w) or d.lookup(w.lower()):
                    continue
                fixes = [v for v in accent_variants(w) if d.lookup(v) or d.lookup(v.lower())]
                if fixes:
                    hits.setdefault((w, fixes[0]), []).append("%s:%d" % (path, n))
for (w, fix), locs in sorted(hits.items(), key=lambda x: -len(x[1])):
    print("%4d  %-22s -> %-22s %s" % (len(locs), w, fix, locs[0]))
print("\n%d formas distintas, %d ocurrencias" % (len(hits), sum(len(v) for v in hits.values())))
