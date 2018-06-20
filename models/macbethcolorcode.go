package models

/*
MacbethColorCode :enum definition for Macbeth color code
*/
type MacbethColorCode int

// DarkSkin :dark skin color
const (
	DarkSkin MacbethColorCode = iota
	LightSkin
	BlueSky
	Foliage
	BlueFlower
	BluishGreen
	Orange
	PurplishBlue
	ModerateRed
	Purple
	YellowGreen
	OrangeYellow
	Blue
	Green
	Red
	Yellow
	Magenta
	Cyan
	White
	Neutral8
	Neutral6p5
	Neutral5
	Neutral3p5
	Black
)

func (ma MacbethColorCode) String() string {
	switch ma {
	case DarkSkin:
		return "DarkSkin"
	case LightSkin:
		return "LightSkin"
	case BlueSky:
		return "BlueSky"
	case Foliage:
		return "Foliage"
	case BlueFlower:
		return "BlueFlower"
	case BluishGreen:
		return "BluishGreen"
	case Orange:
		return "Orange"
	case PurplishBlue:
		return "PurplishBlue"
	case ModerateRed:
		return "ModerateRed"
	case Purple:
		return "Purple"
	case YellowGreen:
		return "YellowGreen"
	case OrangeYellow:
		return "OrangeYellow"
	case Blue:
		return "Blue"
	case Green:
		return "Green"
	case Red:
		return "Red"
	case Yellow:
		return "Yellow"
	case Magenta:
		return "Magenta"
	case Cyan:
		return "Cyan"
	case White:
		return "White"
	case Neutral8:
		return "Neutral8"
	case Neutral6p5:
		return "Neutral6p5"
	case Neutral5:
		return "Neutral5"
	case Neutral3p5:
		return "Neutral3p5"
	case Black:
		return "Black"
	default:
		return ""
	}

}
