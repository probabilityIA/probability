package domain

import "github.com/secamc93/probability/back/central/shared/shippingpkg"

type PackageConfig struct {
	Strategy string
	Boxes    []shippingpkg.Box
}
