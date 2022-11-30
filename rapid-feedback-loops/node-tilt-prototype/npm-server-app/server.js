const express = require('express');
const port = process.env.PORT || 8080;

const app = express();

app.get('/', (request, response) => {
  response.send(`<!DOCTYPE html>
<html>
  <head>
    <title>Powered By Paketo Buildpacks</title>
  </head>
  <body>
    <img style="display: block; margin-left: auto; margin-right: auto; width: 50%;" src="https://paketo.io/images/paketo-logo-full-color.png"></img>
  </body>
</html>`);
});

/* Commmenting out the block above and uncommenting the block below will swap the image on the homepage */
// app.get('/', (request, response) => {
//   response.send(`<!DOCTYPE html>
// <html>
//   <head>
//     <title>Powered By Paketo Buildpacks</title>
//   </head>
//   <body>
//     <img style="display: block; margin-left: auto; margin-right: auto; width: 50%;" src="https://paketo.io/v2/images/buildpack-equation.svg"></img>
//   </body>
// </html>`);
// });
//

app.listen(port);


 // Uncomment the code below ONLY after you have added `chalk` to the project
const chalk = require('chalk');
console.log(chalk.blue("HELLO WORLD!"));
