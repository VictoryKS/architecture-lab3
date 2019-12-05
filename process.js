'use strict';

const fs = require('fs');
const path = require('path');

function correctText(data) {
  const str = data.toString();
  let res = '';
  for (let i = 0; i < str.length; i++) {
    res += str[i];
    if (str[i - 1] && str[i + 1] && str[i] === ','&& !str[i + 1].match(/\s/)
        && (str[i - 1].match(/[a-z]/i) || str[i + 1].match(/[a-z]/i)))
      res += " ";
  }
  return res;
}

function makeDir(path) {
  try {
    if (!fs.existsSync(path)){
      fs.mkdirSync(path);
    }
  } catch (err) {
    throw err;
  }
}

function correction(folderDir, resDir) {
  makeDir(resDir);

  let files = fs.readdirSync(folderDir)
    .filter(fileName => fileName.endsWith('.txt'));

  for (let file of files) {
    fs.readFile(folderDir + '/' + file, 'utf8', (err, data) => {
      if (err) throw err;

      const result = correctText(data);
      if (result) {
        const resFile = resDir + '/' + file.substring(0, file.indexOf('.txt')) + '.res';
        fs.writeFile(resFile, result, (err) => {
          if (err) throw err;
        })
      }
    });
  }
  console.log('Total number of processed files: ' + files.length);
}

// Usage

const folderPath = process.argv[2];
const resultPath = process.argv[3];

if (!folderPath || !resultPath) {
  throw new Error('Incorrect input!');
} else {
  correction(folderPath, resultPath);
}
