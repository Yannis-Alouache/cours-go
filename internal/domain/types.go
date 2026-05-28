package domain

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Room struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Slot struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Available bool   `json:"available"`
}

type Reservation struct {
	ID        int64  `json:"id"`
	RoomID    int64  `json:"room_id"`
	RoomName  string `json:"room_name"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CreateReservationRequest struct {
	RoomID    int64  `json:"room_id"`
	Day       string `json:"day"`
	StartHour int    `json:"start_hour"`
}

type ClientConfig struct {
	APIURL string `json:"api_url"`
	Theme  string `json:"theme"`
}

type ClientState struct {
	Token     string `json:"token"`
	UserEmail string `json:"user_email"`
}
