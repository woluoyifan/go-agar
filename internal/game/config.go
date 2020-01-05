package game

import (
	"github.com/spf13/viper"
	"go-agar/internal/util"
	"time"
)

type config struct {
	Port                    int
	Debug                   bool
	BattlePlayerLimit       int
	TickRate                int
	ScreenWidth             float64
	ScreenHeight            float64
	GameWidth               float64
	GameHeight              float64
	GameMaxMass             float64
	FoodMass                float64
	FoodMaxNum              int
	CellMaxMass             float64
	CellMaxNum              int
	CellDefaultSpeed        float64
	CellMergeDistanceRate   float64
	VirusMinMass            float64
	VirusMaxMass            float64
	VirusMaxNum             int
	VirusColor              string
	MinMassLose             float64
	MassLoseRate            float64
	MaxHeartbeatInterval    time.Duration
	DefaultPlayerMass       float64
	SlowBase                float64
	InitMassLog             float64
	MergeInterval           time.Duration
	MassWinRate             float64
	FireFoodMass            float64
	FireFoodRate            float64
	FireFoodSpeedSlow       float64
	FireFoodSpeed           float64
	SplitSpeed              float64
	AnonymousUserNamePrefix string
}

var (
	Config *config
)

func init() {
	setDefaultConfig()
	viper.SetConfigName("game")
	viper.AddConfigPath("./")
	viper.ReadInConfig()
	Config = &config{
		Port:                    viper.GetInt("Port"),
		Debug:                   viper.GetBool("Debug"),
		BattlePlayerLimit:       viper.GetInt("BattlePlayerLimit"),
		TickRate:                viper.GetInt("TickRate"),
		ScreenWidth:             float64(viper.GetInt("ScreenWidth")),
		ScreenHeight:            float64(viper.GetInt("ScreenHeight")),
		GameWidth:               float64(viper.GetInt("GameWidth")),
		GameHeight:              float64(viper.GetInt("GameHeight")),
		GameMaxMass:             viper.GetFloat64("GameMaxMass"),
		FoodMass:                viper.GetFloat64("FoodMass"),
		FoodMaxNum:              viper.GetInt("FoodMaxNum"),
		CellMaxMass:             viper.GetFloat64("CellMaxMass"),
		CellMaxNum:              viper.GetInt("CellMaxNum"),
		CellDefaultSpeed:        viper.GetFloat64("CellDefaultSpeed"),
		CellMergeDistanceRate:   viper.GetFloat64("CellMergeDistanceRate"),
		VirusMinMass:            viper.GetFloat64("VirusMinMass"),
		VirusMaxMass:            viper.GetFloat64("VirusMaxMass"),
		VirusMaxNum:             viper.GetInt("VirusMaxNum"),
		VirusColor:              viper.GetString("VirusColor"),
		MinMassLose:             viper.GetFloat64("MinMassLose"),
		MassLoseRate:            viper.GetFloat64("MassLoseRate"),
		MaxHeartbeatInterval:    viper.GetDuration("MaxHeartbeatInterval"),
		DefaultPlayerMass:       viper.GetFloat64("DefaultPlayerMass"),
		SlowBase:                viper.GetFloat64("SlowBase"),
		InitMassLog:             viper.GetFloat64("InitMassLog"),
		MergeInterval:           viper.GetDuration("MergeInterval"),
		MassWinRate:             viper.GetFloat64("MassWinRate"),
		FireFoodMass:            viper.GetFloat64("FireFoodMass"),
		FireFoodRate:            viper.GetFloat64("FireFoodRate"),
		FireFoodSpeedSlow:       viper.GetFloat64("FireFoodSpeedSlow"),
		FireFoodSpeed:           viper.GetFloat64("FireFoodSpeed"),
		SplitSpeed:              viper.GetFloat64("SplitSpeed"),
		AnonymousUserNamePrefix: viper.GetString("AnonymousUserNamePrefix"),
	}
}

func setDefaultConfig() {
	defaultPlayerMass := float64(10)
	slowBase := 4.5
	viper.SetDefault("Port", 38888)
	viper.SetDefault("Debug", false)
	viper.SetDefault("BattlePlayerLimit", 50)
	viper.SetDefault("TickRate", 60)
	viper.SetDefault("ScreenWidth", 1024)
	viper.SetDefault("ScreenHeight", 768)
	viper.SetDefault("GameWidth", 5000)
	viper.SetDefault("GameHeight", 5000)
	viper.SetDefault("GameMaxMass", 20000)
	viper.SetDefault("FoodMass", 1)
	viper.SetDefault("FoodMaxNum", 1000)
	viper.SetDefault("CellMaxMass", 100)
	viper.SetDefault("CellMaxNum", 16)
	viper.SetDefault("CellDefaultSpeed", 6.25)
	viper.SetDefault("CellMergeDistanceRate", 1.75)
	viper.SetDefault("VirusMinMass", 100)
	viper.SetDefault("VirusMaxMass", 150)
	viper.SetDefault("VirusMaxNum", 50)
	viper.SetDefault("VirusColor", "#7bff66")
	viper.SetDefault("MinMassLose", 50)
	viper.SetDefault("MassLoseRate", 1)
	viper.SetDefault("MaxHeartbeatInterval", 5*time.Second)
	viper.SetDefault("DefaultPlayerMass", defaultPlayerMass)
	viper.SetDefault("SlowBase", slowBase)
	viper.SetDefault("InitMassLog", util.Log(defaultPlayerMass, slowBase))
	viper.SetDefault("MergeInterval", 15*time.Second)
	viper.SetDefault("MassWinRate", 1.1)
	viper.SetDefault("FireFoodMass", 20)
	viper.SetDefault("FireFoodRate", 0)
	viper.SetDefault("FireFoodSpeedSlow", 0.5)
	viper.SetDefault("FireFoodSpeed", 25)
	viper.SetDefault("SplitSpeed", 25)
	viper.SetDefault("AnonymousUserNamePrefix", "u")
}
