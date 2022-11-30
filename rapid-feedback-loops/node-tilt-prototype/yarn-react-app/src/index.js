import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import Question from './question';

ReactDOM.render(
  <React.StrictMode>
    <App />
  //Uncomment the line below to see changes on the homepage
    // <Question />
  </React.StrictMode>,
  document.getElementById('root')
);
