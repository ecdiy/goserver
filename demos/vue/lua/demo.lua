id=db:QueryString("select Id from GkUser where Username='test'")

return {
    abc=param("abc"),
    Id =id,
    A = 10
}
