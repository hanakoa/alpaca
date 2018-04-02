(function() {
  // We use an IIFE so we can declare variables
  // (like these) without polluting global scope.
  const scheme = 'http';
  const baseUrl = 'api.alpaca.minikube';
  const configObjectName = 'alpacaConfig';

  if (!window[configObjectName]) {
    window[configObjectName] = {};
  }

  const apiNames = ['auth', 'passwordReset', 'help', 'accountRequests'];

  // Improved loop, assigns length to variable & defines i to avoid hoisting
  let dataLength = apiNames.length,
    i;
  for (i = 0; i < dataLength; i++) {
    let apiName = apiNames[i];
    let fullApiUrl = `${scheme}://${baseUrl}/${apiName}`;

    if (!window[configObjectName][apiName]) {
      window[configObjectName][apiName] = {};
    }

    window[configObjectName][apiName].baseApiUrl = fullApiUrl;
  }
})();
