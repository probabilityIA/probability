package usecaseshipment

import (
	"bytes"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var carrierLogoURLs = map[string]string{
	"SERVIENTREGA":                 "https://i.revistapym.com.co/old/2021/09/WhatsApp-Image-2021-09-25-at-1.08.55-PM.jpeg?w=400&r=1_1",
	"COORDINADORA":                 "https://olartemoure.com/wp-content/uploads/2023/05/coordinadora-logo.png",
	"DHL":                          "https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png",
	"DHLEXPRESS":                   "https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png",
	"FEDEX":                        "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/FedEx_Express.svg/960px-FedEx_Express.svg.png",
	"INTERRAPIDISIMO":              "https://interrapidisimo.com/wp-content/uploads/Logo-Inter-Rapidisimo-Vv-400x431-1.png",
	"472LOGISTICA":                 "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s",
	"ENVIA":                        "https://images.seeklogo.com/logo-png/31/1/envia-mensajeria-logo-png_seeklogo-311137.png",
	"TCC":                          "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png",
	"TRANSPORTADORADECARACOLOMBIA": "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png",
	"DEPRISA":                      "https://www.specialcolombia.com/wp-content/uploads/2023/05/Logo_azul_concepto_azul-deprisa.png",
	"ENVIOCLICK":                   "https://www.envioclickpro.com.co/assets/images/envioclick-logo.png",
	"MIPAQUETE":                    "https://mipaquete.com/wp-content/uploads/2021/03/mipaquete-logo.png",
}

var (
	carrierLogoCache   = map[string][]byte{}
	carrierLogoCacheMu sync.RWMutex
	carrierLogoClient  = &http.Client{Timeout: 6 * time.Second}
)

func carrierLogoKey(carrier string) string {
	k := strings.ToUpper(carrier)
	repl := strings.NewReplacer(" ", "", "-", "", "_", "")
	return repl.Replace(k)
}

func getCarrierLogoBytes(carrier string) []byte {
	key := carrierLogoKey(carrier)
	if key == "" {
		return nil
	}
	url, ok := carrierLogoURLs[key]
	if !ok {
		return nil
	}

	carrierLogoCacheMu.RLock()
	if data, hit := carrierLogoCache[key]; hit {
		carrierLogoCacheMu.RUnlock()
		return data
	}
	carrierLogoCacheMu.RUnlock()

	resp, err := carrierLogoClient.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	data := buf.Bytes()

	carrierLogoCacheMu.Lock()
	carrierLogoCache[key] = data
	carrierLogoCacheMu.Unlock()
	return data
}
