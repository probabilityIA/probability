#!/usr/bin/env python3
"""Revisa la ortografia del espanol en los textos de cara al usuario del front."""
import argparse
import os
import re
import sys

SKILL_DIR = os.path.dirname(os.path.abspath(__file__))
EXTS = {".tsx", ".ts", ".jsx", ".js", ".astro", ".mdx"}
SKIP_DIRS = {"node_modules", ".next", "dist", "build", ".git", "coverage", ".turbo", "public"}
DEFAULT_PATHS = ["front/central/src", "front/website/src"]
ASCII_ONLY_LINES = 500

LETTER = "A-Za-z0-9_À-ſ"
ESCAPES = {
    "á": "\\u00e1", "é": "\\u00e9", "í": "\\u00ed",
    "ó": "\\u00f3", "ú": "\\u00fa", "ñ": "\\u00f1",
    "ü": "\\u00fc", "Á": "\\u00c1", "É": "\\u00c9",
    "Í": "\\u00cd", "Ó": "\\u00d3", "Ú": "\\u00da",
    "Ñ": "\\u00d1", "Ü": "\\u00dc", "¿": "\\u00bf",
    "¡": "\\u00a1",
}
ATTR_SKIP = re.compile(
    r"(?:className|class|href|src|srcSet|id|key|name|type|path|route|to|rel|target|"
    r"htmlFor|testID|data-[\w-]+|aria-[\w-]+|role|accept|pattern|format|locale|"
    r"icon|variant|color|size|as|method|action|fill|stroke|style|transform|"
    r"points|viewBox|d|stopColor|offset|placeholder-\w+)\s*=\s*[\{\(\[]?\s*$"
)
LINE_SKIP = re.compile(r"^\s*(?:import|export\s+(?:\*|\{)|//|/\*|\*)|require\s*\(|\bfrom\s+['\"]")
PATHY = re.compile(r"^[\w./@:#?&=%~+-]+$")
CODEY = re.compile(r"^[a-z][a-z0-9]*(?:[_-][a-z0-9]+)*$")
WORDY = re.compile(r"[A-Za-zÀ-ſ]+")
JSX_CODEY = re.compile(r"===|!==|&&|\|\||=>|\?\?|^\s*[:?]|\w\(|^\s*\d+\s")


def load_pairs(path, phrase):
    out = []
    with open(path, encoding="utf-8") as fh:
        for raw in fh:
            line = raw.rstrip("\n")
            if not line.strip() or line.lstrip().startswith("#"):
                continue
            parts = line.split("\t")
            if len(parts) < 2:
                continue
            bad, good = parts[0].strip(), parts[1].strip()
            mode = parts[2].strip() if len(parts) > 2 and parts[2].strip() else "auto"
            out.append((bad, good, mode, phrase))
    return out


def case_variants(bad, good):
    yield bad, good
    if bad[:1].islower():
        yield bad[:1].upper() + bad[1:], good[:1].upper() + good[1:]
    if bad.lower() == bad:
        yield bad.upper(), good.upper()


def build_matcher(pairs, with_manual):
    table = {}
    for bad, good, mode, phrase in pairs:
        if mode == "manual" and not with_manual:
            continue
        for b, g in case_variants(bad, good):
            table.setdefault(b, (g, mode, phrase))
    keys = sorted(table, key=len, reverse=True)
    if not keys:
        return None, table
    alt = "|".join(re.escape(k) for k in keys)
    rx = re.compile(r"(?<![%s])(?:%s)(?![%s])" % (LETTER, alt, LETTER))
    return rx, table


SPAN_RX = [
    ("s", re.compile(r"'((?:[^'\\\n]|\\.)*)'")),
    ("d", re.compile(r"\"((?:[^\"\\\n]|\\.)*)\"")),
    ("t", re.compile(r"`((?:[^`\\]|\\.)*)`")),
    ("jsx", re.compile(r">([^<>{}]+)<")),
]


def spans(line):
    taken = []
    for kind, rx in SPAN_RX:
        for m in rx.finditer(line):
            a, b = m.start(1), m.end(1)
            if any(a < tb and ta < b for ta, tb in taken):
                continue
            text = m.group(1)
            if not text.strip():
                continue
            if ATTR_SKIP.search(line[max(0, m.start() - 60):m.start()]):
                continue
            if " " not in text:
                bare = text.strip(":*!?.,;-")
                if not WORDY.fullmatch(bare) or bare.islower():
                    continue
            if kind == "jsx" and JSX_CODEY.search(text):
                continue
            taken.append((a, b))
            yield kind, a, text


APERTURA = re.compile(r"[.!?]\s+|^")


def scan_file(path, rx, table, want_apertura):
    try:
        src = open(path, encoding="utf-8").read()
    except (UnicodeDecodeError, OSError):
        return [], 0
    lines = src.split("\n")
    findings = []
    for n, line in enumerate(lines, 1):
        if LINE_SKIP.search(line):
            continue
        for kind, off, text in spans(line):
            if rx is not None:
                for m in rx.finditer(text):
                    good, mode, phrase = table[m.group(0)]
                    findings.append({
                        "line": n, "col": off + m.start() + 1, "kind": kind,
                        "bad": m.group(0), "good": good, "mode": mode,
                        "phrase": phrase, "span_start": off, "in_span": m.start(),
                        "type": "frase" if phrase else "tilde", "text": text,
                    })
            if want_apertura and "${" not in text:
                stripped = text.strip()
                if len(stripped) > 6 and " " in stripped:
                    if stripped.endswith("?") and "¿" not in stripped:
                        findings.append({"line": n, "col": off + 1, "kind": kind,
                                         "bad": stripped[:48], "good": "falta ¿ de apertura",
                                         "mode": "manual", "type": "apertura", "text": text})
                    if stripped.endswith("!") and "¡" not in stripped and not stripped.endswith("!!"):
                        findings.append({"line": n, "col": off + 1, "kind": kind,
                                         "bad": stripped[:48], "good": "falta ¡ de apertura",
                                         "mode": "manual", "type": "apertura", "text": text})
    return findings, len(lines)


def to_ascii(word):
    return "".join(ESCAPES.get(ch, ch) for ch in word)


def render(good, kind, ext, ascii_mode):
    if not ascii_mode:
        return good
    if kind == "jsx":
        if ext == ".astro":
            return "".join("&#x%x;" % ord(c) if ord(c) > 127 else c for c in good)
        return "{'%s'}" % to_ascii(good)
    return to_ascii(good)


def fix_file(path, findings, total_lines):
    ext = os.path.splitext(path)[1]
    ascii_mode = total_lines >= ASCII_ONLY_LINES
    by_line = {}
    for f in findings:
        if f["type"] == "apertura" or f["mode"] != "auto":
            continue
        by_line.setdefault(f["line"], []).append(f)
    if not by_line:
        return 0
    lines = open(path, encoding="utf-8").read().split("\n")
    changed = 0
    for n, items in by_line.items():
        line = lines[n - 1]
        for f in sorted(items, key=lambda x: x["col"], reverse=True):
            start = f["col"] - 1
            end = start + len(f["bad"])
            if line[start:end] != f["bad"]:
                continue
            line = line[:start] + render(f["good"], f["kind"], ext, ascii_mode) + line[end:]
            changed += 1
        lines[n - 1] = line
    if changed:
        open(path, "w", encoding="utf-8").write("\n".join(lines))
    return changed


def walk(paths):
    for p in paths:
        if os.path.isfile(p):
            if os.path.splitext(p)[1] in EXTS:
                yield p
            continue
        for root, dirs, files in os.walk(p):
            dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
            for name in sorted(files):
                if os.path.splitext(name)[1] in EXTS:
                    yield os.path.join(root, name)


def main():
    ap = argparse.ArgumentParser(description="Ortografia del espanol en textos del front")
    ap.add_argument("paths", nargs="*", default=None)
    ap.add_argument("--fix", action="store_true", help="aplica las correcciones marcadas auto")
    ap.add_argument("--ambiguos", action="store_true", help="incluye palabras de tilde diacritica")
    ap.add_argument("--sin-apertura", action="store_true", help="no revisar signos ? y !")
    ap.add_argument("--limite", type=int, default=60, help="maximo de hallazgos mostrados por tipo")
    args = ap.parse_args()

    paths = args.paths or [p for p in DEFAULT_PATHS if os.path.exists(p)]
    if not paths:
        print("No hay rutas para revisar. Ejecutalo desde la raiz del repo.", file=sys.stderr)
        return 2

    pairs = load_pairs(os.path.join(SKILL_DIR, "diccionario.tsv"), False)
    pairs += load_pairs(os.path.join(SKILL_DIR, "frases.tsv"), True)
    rx, table = build_matcher(pairs, args.ambiguos)

    total, fixed, shown = 0, 0, 0
    per_file = []
    for path in walk(paths):
        findings, nlines = scan_file(path, rx, table, not args.sin_apertura)
        if not findings:
            continue
        total += len(findings)
        per_file.append((path, findings, nlines))

    per_file.sort(key=lambda x: -len(x[1]))
    for path, findings, nlines in per_file:
        if args.fix:
            fixed += fix_file(path, findings, nlines)
            continue
        if shown >= args.limite:
            continue
        marca = " [ASCII]" if nlines >= ASCII_ONLY_LINES else ""
        print("\n%s  (%d lineas%s)  %d hallazgos" % (path, nlines, marca, len(findings)))
        for f in findings[:12]:
            tag = {"tilde": "tilde", "frase": "frase", "apertura": "signo"}[f["type"]]
            auto = "auto" if f["mode"] == "auto" and f["type"] != "apertura" else "rev "
            print("  %5d:%-4d %s %s  %s -> %s" % (f["line"], f["col"], auto, tag, f["bad"], f["good"]))
        if len(findings) > 12:
            print("  ... %d mas" % (len(findings) - 12))
        shown += 1

    print("\n%d archivo(s) con hallazgos, %d hallazgo(s) en total." % (len(per_file), total))
    if args.fix:
        print("%d correccion(es) aplicada(s). Revisa el diff antes de commitear." % fixed)
    elif shown < len(per_file):
        print("Mostrados %d archivos de %d. Sube --limite para ver mas." % (shown, len(per_file)))
    if args.fix:
        return 0
    return 1 if total else 0


if __name__ == "__main__":
    sys.exit(main())
