package route

type (
	Upsert struct {
		ID int64 `json:"route_id"`
		Insert
	}

	Insert struct {
		Name  string  `json:"route_name"`
		Load  float64 `json:"load"`
		Cargo string  `json:"cargo_type"`
	}

	Get struct {
		ID int64 `params:"id"`
	}

	Delete struct {
		Ids []int64 `json:"route_ids"`
	}

	CreateResp struct {
		ID int64 `json:"route_id"`
	}

	Route struct {
		Upsert
		Actual bool `json:"-"`
	}
)
