package models

/*
SettingInfo :strcutrue for setting information
*/
type SettingInfo struct {
	Gamma       float64 // gamma value
	LightSource string  // light source info
	PatchSize   struct {
		H int
		V int
	}
	StdPatchSavePath     string // file save path of standard Macbeth chart
	StdPatchSaveDirName  string
	DevPatchSavePath     string // file save path of device Macbeth chart
	DevPatchSaveDirName  string
	DeltaESavePath       string // file save path of delta E
	DeltaESaveDirName    string
	DeiceQEDataPath      string // QE data path
	WhitePixelDataPath   string // White pixel data path
	LinearMatrixDataPath string // Linear matrix data path
}
