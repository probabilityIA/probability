package entities

import "testing"

func TestResolveSiigoBaseURL(t *testing.T) {
	const prod = "https://api.siigo.com"
	const test = "https://api-sandbox.siigo.com"
	const override = "https://custom.siigo.example"

	cases := []struct {
		name      string
		isTesting bool
		baseTest  string
		override  string
		base      string
		want      string
	}{
		{
			name: "produccion sin api_url: cae al base_url del integration_type (bug de los webhooks)",
			base: prod,
			want: prod,
		},
		{
			name:      "modo pruebas: gana base_url_test",
			isTesting: true,
			baseTest:  test,
			override:  override,
			base:      prod,
			want:      test,
		},
		{
			name:      "modo pruebas sin base_url_test: usa el override",
			isTesting: true,
			override:  override,
			base:      prod,
			want:      override,
		},
		{
			name:     "el override api_url le gana a base_url",
			override: override,
			base:     prod,
			want:     override,
		},
		{
			name:     "espacios en blanco no cuentan como URL",
			override: "   ",
			base:     prod,
			want:     prod,
		},
		{
			name: "sin ninguna URL configurada devuelve vacio (el llamador corta con error claro)",
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ResolveSiigoBaseURL(tc.isTesting, tc.baseTest, tc.override, tc.base)
			if got != tc.want {
				t.Fatalf("ResolveSiigoBaseURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
