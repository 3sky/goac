package api


//HelloStruct for valid JSON
type HelloStruct struct {
	ID  int    `json:"ID,omitempty"`
	SAY string `json:"INFO,omitempty"`
}

// Similar struct to db.StatusStruct, but without model
type AppStatusStruct struct {
	ID uint 
	APP_NAME string `json:"app_name"` 
	APP_VERSION string `json:"app_version"`
	UPDATE_DATE string`json:"updated_date"`
	UPDATE_BY string `json:"updated_by"`
}
