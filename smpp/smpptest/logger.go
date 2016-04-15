package smpptest

import (
	"fmt"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	infoHandle := os.Stdout
	warningHandle := os.Stdout
	errorHandle := os.Stdout

	Info = log.New(infoHandle,
		fmt.Sprintf("%s INFO: ", DefaultSystemID),
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		fmt.Sprintf("%s WARNING: ", DefaultSystemID),
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		fmt.Sprintf("%s ERROR: ", DefaultSystemID),
		log.Ldate|log.Ltime|log.Lshortfile)
}
