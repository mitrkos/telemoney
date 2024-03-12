package myerror

type MyError struct {
	Msg string
}


func (e *MyError) Error() string { 
    return e.Msg
}
