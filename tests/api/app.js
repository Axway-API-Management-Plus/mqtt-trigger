var express = require('express')
var app = express()
var bodyParser = require('body-parser')

process.on('SIGTERM', function () {
  process.exit(0);
});

app.use(bodyParser.json())

app.all(/\/.*/, function (req, res) {
  console.log(req.method, req.path)
  console.log("HEADERS:", req.headers)
  console.log("BODY:", req.body)
  res.send({ })
})

app.listen(3000, function () {
  console.log('Example app listening on port 3000!')
})
