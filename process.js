'use strict';

const fs = require('fs');
const path = require('path');

function correctText(str, prevEnd) {
  let res = '';
  if (prevEnd && prevEnd[1] === ',' && !str[0].match(/\s/) &&
     (prevEnd[0].match(/[a-z]/i) || str[0].match(/[a-z]/i)))
     res += " ";
  for (let i = 0; i < str.length - 1 ; i++) {
    res += str[i];
    if (str[i] === ',' && !str[i + 1].match(/\s/)
        && (((str[i - 1] && str[i - 1].match(/[a-z]/i))
        || (prevEnd && prevEnd[1].match(/[a-z]/i)))
        || str[i + 1].match(/[a-z]/i)))
      res += " ";
  }
  res += str[str.length - 1];
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

  let counter = 0;
  for (let file of files) {
    let prevDataEnd;
    const stream = fs.createReadStream(folderDir + '/' + file, { highWaterMark: 64 })
      .on('data', (data) => {
        data = data.toString();
        const result = correctText(data, prevDataEnd);
        const resFile = resDir + '/' + file.substring(0, file.indexOf('.txt')) + '.res';
        if (!prevDataEnd)
          fs.writeFile(resFile, result, (err) => {
            if (err) throw err;
          });
        else
          fs.appendFile(resFile, result, (err) => {
            if (err) throw err;
          });
        prevDataEnd = data.substring(data.length - 2, data.length);
      })
      .on('end', () => {
        counter++;
        if (counter == files.length)
          console.log('Total number of processed files: ' + files.length);
      });
  }
}

// Usage

const folderPath = process.argv[2];
const resultPath = process.argv[3];

if (!folderPath || !resultPath) {
  throw new Error('Incorrect input!');
} else {
  correction(folderPath, resultPath);
}
