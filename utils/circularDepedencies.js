// This function returns a replacer function that can be used with JSON.stringify to prevent circular references from being included in the output string.

function getCircularReplacer() {
  const seen = new WeakSet();
  return (key, value) => {
    try {
      /*
      If the value is an object and not null:
        If we've seen the value before:
          Return undefined to skip this value.
        Otherwise:
          Add the value to the set of seen values.
      */
      if (typeof value === "object" && value !== null) {
        if (seen.has(value)) {
          return;
        }
        seen.add(value);
      }
      return value;
    } catch (err) {
      // Ignore error. Just return undefined.
    }
  };
}

module.exports = getCircularReplacer;
