package logger

import "go.uber.org/zap"

	var Log *zap.Logger

	func Logger(){
		var err error
		Log, err = zap.NewProduction()
		if err != nil {
			panic("Failled to initialize logger:" + err.Error())
		}
	}