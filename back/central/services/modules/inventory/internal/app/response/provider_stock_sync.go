package response

type ProviderSyncResult struct {
	Total     int
	Updated   int
	Unchanged int
	Skipped   int
	Failed    int
}
