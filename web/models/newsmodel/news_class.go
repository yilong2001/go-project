package newsmodel

func GetNewsClassCode(name string) int {
	switch name {
	case "互金":
		return 100100
	case "银行":
		return 110100
	case "保险":
		return 120100
	case "证券":
		return 130100
	case "投资":
		return 140100
	case "信托":
		return 150100
	case "财税":
		return 160100
	case "经济":
		return 180100
	case "政策":
		return 260100
	case "证券执业":
		return 310100
	case "银行执业":
		return 320100
	case "保险执业":
		return 330100
	case "金融执业":
		return 340100
	case "基金执业":
		return 350100
	case "财税执业":
		return 360100
	case "金融工具":
		return 410100
	case "软件技能":
		return 450100
	}

	return 0
}
