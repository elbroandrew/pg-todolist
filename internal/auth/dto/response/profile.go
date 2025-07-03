package response

type Profile struct {
    ID        uint   `json:"id"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}