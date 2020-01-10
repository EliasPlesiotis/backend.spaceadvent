package kuiper

import "strings"

// Responde to user from Kuiper the chatbot
func Responde(text string) string {
	var resp string
	if strings.Contains(strings.ToUpper(text), "HELLO") {
		resp += "Hi there"
	} else if strings.Contains(strings.ToUpper(text), "MOVE") {
		resp += "Your spaceship follows the mouse pointer."
	} else if strings.Contains(strings.ToUpper(text), "BULLETS") {
		resp += "Gather bullets in order to shoot"
	} else if strings.Contains(strings.ToUpper(text), "SHOOT") {
		resp += "To shoot click with the mouse and a bullet will spawn with direction from the spaceship to the mouse"
	} else if strings.Contains(strings.ToUpper(text), "SHIELD") {
		resp += "Gather power from the quantom fields and use the space bar to use it for shield"
	} else if strings.Contains(strings.ToUpper(text), "GRAVITY") {
		resp += "Neutron stars and Black holes will pull your spacecraft with their gravity, so be carefull"
	} else if strings.Contains(strings.ToUpper(text), "THANKS") {
		resp += "You are welcome"
	} else {
		resp += "I dont understand the question"
	}

	return resp
}
