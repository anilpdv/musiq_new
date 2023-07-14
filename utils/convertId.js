const test = length => {
  let possible;
  let randomChar = "";

  possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

  for (let i = 0; i < length; i++) {
    randomChar =
      randomChar + possible.charAt(Math.floor(Math.random() * possible.length));
  }
  return randomChar;
};

const convertId = videoId => {
  let charid;
  let length;
  let generatedValue = "";

  for (let i = 0; i < videoId.length; i++) {
    charid = videoId.charCodeAt(i) + 13;

    length = Math.floor(Math.random() * 3 + 1);

    generatedValue = generatedValue + (charid + test(length));
  }
  return generatedValue;
};

module.exports = convertId;
