package app

import "github.com/zgs225/go-ecdict/dict"

func initHookECDict(app *Application) error {
	dict, err := dict.NewSimpleDict(app.Config.GetString("ecdict"))
	if err != nil {
		return err
	}
	app.Dict = dict
	app.Logger.Info("ECDICT 初始化完成")
	return nil
}
