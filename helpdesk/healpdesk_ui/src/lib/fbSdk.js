export const initFacebookSdk = () => {
    return new Promise((resolve, _reject) => {
      // Load the Facebook SDK asynchronously
      window.fbAsyncInit = () => {
        console.log('Loading Facebook SDK', window, window.FB)

        window.FB.init({
          appId: '445200227996341',
          cookie: true,
          xfbml: true,
          version: 'v19.0'
        })
        // Resolve the promise when the SDK is loaded
        resolve()
      }
    })
}

export const getFacebookLoginStatus = () => {
    return new Promise((resolve, reject) => {
      window.FB.getLoginStatus((response) => {
        resolve(response);
      });
    });
  };

  export const fbLogin = () => {
    return new Promise((resolve, reject) => {
      window.FB.login((response) => {
        resolve(response)        
      })
    })
  }

