const http = require('http')
const port = process.env.PORT || 8080

const requestHandler = (request, response) => {
    response.end("hello world")
}

// Commmenting out the block above and uncommenting the block below will change the message on the homepage.
//
// const requestHandler = (request, response) => {
    // response.end("goodbye world")
// }

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
    if (err) {
        return console.log('something bad happened', err)
    }

    console.log(`server is listening on ${port}`)
})
