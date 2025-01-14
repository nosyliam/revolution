package hive

import (
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

var ClaimHiveImage = ImageSteps{
	SelectY1(Change, Add(S[int]("honeyOffsetY"), 23)),
	SelectY2(Change, Add(Y1, 30)),
	Variance(0),
	Direction(0),
	Search("claimhive").Find(),
}

var SendTradeImage = ImageSteps{
	SelectY1(Change, Add(S[int]("honeyOffsetY"), 23)),
	SelectY2(Change, Add(Y1, 30)),
	Variance(0),
	Direction(0),
	Search("sendtrade").Find(),
}

var TradeDisabledImage = ImageSteps{
	SelectY1(Change, Add(S[int]("honeyOffsetY"), 23)),
	SelectY2(Change, Add(Y1, 30)),
	Variance(0),
	Direction(0),
	Search("tradedisabled").Find(),
}

var TradeLockedImage = ImageSteps{
	SelectY1(Change, Add(S[int]("honeyOffsetY"), 23)),
	SelectY2(Change, Add(Y1, 30)),
	Variance(0),
	Direction(0),
	Search("tradelocked").Find(),
}

var AllHiveImages = ImageSteps{
	SelectY1(Change, Add(S[int]("honeyOffsetY"), 23)),
	SelectY2(Change, Add(Y1, 30)),
	Variance(0),
	Direction(0),
	Search("claimhive", "sendtrade", "tradelocked", "tradedisabled").Find(),
}
