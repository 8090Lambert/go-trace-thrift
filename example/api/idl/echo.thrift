namespace go echo

service Echo{
    string  Emit(1:string message)
}