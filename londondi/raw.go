package londondi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/byuoitav/london-audio-microservice/connect"
	"github.com/fatih/color"
)

func HandleRawCommandString(rawCommand RawDICommand) error {

	log.Printf("Handling raw command: %s...", rawCommand.Command)

	connection, connectError := net.Dial("tcp", rawCommand.Address+":"+
		rawCommand.Port)
	if connectError != nil {
		log.Printf(connectError.Error())
		return connectError
	}

	log.Printf("Converting to command to hex value...")
	hexCommand, hexError := hex.DecodeString(rawCommand.Command)
	if hexError != nil {
		fmt.Println(hexError.Error())
		return hexError
	}

	log.Printf("hexCommand: %x", hexCommand)

	_, writeError := connection.Write(hexCommand)
	if writeError != nil {
		log.Printf(writeError.Error())
	}

	connection.Close()
	return connectError
}

func HandleRawCommandBytes(command []byte, address string) error {

	log.Printf("%s", color.HiCyanString("Handling raw command %x with address %s...", command, address))

	connection, err := connect.GetConnection(address)
	if err != nil {
		msg := fmt.Sprintf("problem getting connection to device: %s", err.Error())
		return errors.New(msg)
	}

	_, err = connection.Write(command)
	if err != nil {

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {

			_, err = connect.HandleTimeout(connection, command, connect.Write)
		}
	}

	if err != nil {

		msg := color.HiRedString("could not write to device: ", err.Error())
		log.Printf("%s %s", color.HiRedString("[error]"), msg)
		return errors.New(msg)
	}

	return nil
}
