package londondi

import (
	"bufio"
	"encoding/hex"
	"errors"
	"log"
	"net"

	"github.com/byuoitav/av-api/status"
)

func GetVolume(address, input string) (status.Volume, error) {

	command, err := BuildSubscribeCommand(address, input, "volume", DI_SUBSCRIBESVPERCENT)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	state, err := ParseVolumeStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	return state, nil

}

func GetMute(address, input string) (status.MuteStatus, error) {

	command, err := BuildSubscribeCommand(address, input, "mute", DI_SUBSCRIBESV)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	state, err := ParseMuteStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	return state, nil

}

//@param subscribe - TRUE indicates subscribe, FALSE indicates unsubscribe
func BuildSubscribeCommand(address, input, state string, messageType int32) (RawDICommand, error) {

	log.Printf("Building raw command to subsribe to %s of input %s on address %s", state, input, address)

	command := []byte{byte(messageType)}

	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, NODE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stateVariable, _ := hex.DecodeString(stateVariables[state])
	command = append(command, stateVariable...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, RATE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	checksum := GetChecksumByte(command)
	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	STX := []byte{byte(reserved["STX"])}
	command = append(STX, command...)
	ETX := byte(reserved["ETX"])
	command = append(command, ETX)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	return RawDICommand{
		Address: address,
		Port:    PORT,
		Command: hex.EncodeToString(command),
	}, nil
}

func HandleStatusCommand(subscribe RawDICommand) ([]byte, error) {

	log.Printf("Handling status command...")

	connection, err := net.Dial("tcp", subscribe.Address+":"+subscribe.Port)
	if err != nil {
		errorMessage := "Could not connect to device: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	defer connection.Close()

	log.Printf("Converting command to hex value...")
	command, err := hex.DecodeString(subscribe.Command)
	if err != nil {
		errorMessage := "Could not convert command to hex value: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	_, err = connection.Write(command)
	if err != nil {
		errorMessage := "Could not send message to device: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	//	var reader bufio.Reader
	//	response, _ reader.ReadBytes(ACK)

	reader := bufio.NewReader(connection)

	response, err := reader.ReadBytes(ETX)
	if err != nil {
		errorMessage := "Could not find ETX: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	log.Printf("Message: %x", response)

	return response, nil

}

func ParseVolumeStatus(message []byte) (status.Volume, error) {

	return status.Volume{}, nil
}

func ParseMuteStatus(message []byte) (status.MuteStatus, error) {

	return status.MuteStatus{}, nil
}