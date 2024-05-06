import React, { useEffect, useState } from "react";
import Card from "./CommonComponents/Card";
import Button from "./CommonComponents/Button";
import { useAuth } from "../hooks/auth";
import FacebookLogin from "react-facebook-login/dist/facebook-login-render-props";
import { GraphApi } from "../Api/Axios";
import { showError, showSuccess } from "../lib/utils";
import { Link } from "react-router-dom";
import { fbLogin, getFacebookLoginStatus, initFacebookSdk } from "../lib/fbSdk";
const FBConnectPage = () => {
  const [loading, setLoading] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const auth = useAuth();

  const facebookAppID = process.env.REACT_APP_FACEBOOK_APP_ID || 445200227996341;
  const facebookRedirectURI = process.env.REACT_APP_PUBLIC_URL_ENCODED || "https://e665-2405-201-d015-c112-44af-904d-4861-c76b.ngrok-free.app/connect-page"

  const getPageId = async (accessToken) => {
    setLoading(true);
    try {
      console.log('graphapi', GraphApi)
      const res = await GraphApi.get("/me/accounts", {
        params: { access_token: accessToken },
      });
      const pageObj = {
        name: res?.data?.data[0]?.name,
        id: res?.data?.data[0]?.id,
        pageAccessToken: res?.data?.data[0]?.access_token,
      };

      // const pageObj = {
      //   name: "Umed",
      //   id: "12334",
      //   pageAccessToken: accessToken,
      // };

      localStorage.setItem("FB_PAGE_ID", pageObj?.id);
      localStorage.setItem("FB_PAGE_ACCESS_TOKEN", pageObj.pageAccessToken);
      localStorage.setItem("FB_PAGE_DETAILS", JSON.stringify(pageObj));
      showSuccess(`Connected to Facebook page ${pageObj?.name}`);
      setIsConnected(true);
    } catch (error) {
      setLoading(false);
      showError("Could not connect to the Facebook page");
      localStorage.removeItem("FB_ACCESS_TOKEN");
      localStorage.removeItem("FB_PAGE_ACCESS_TOKEN");
      localStorage.removeItem("FB_PAGE_ID");
    }
    setLoading(false);
  };

  const responseFacebook = async (data) => {
    setLoading(true);
    try {
      if (data.accessToken) {
        const accessToken = data.accessToken;
        localStorage.setItem("FB_ACCESS_TOKEN", accessToken);
        await getPageId(accessToken);
      }
    } catch (error) {
      setLoading(false);
      showError("Could not connect to the Facebook page");
      localStorage.removeItem("FB_ACCESS_TOKEN");
      localStorage.removeItem("FB_PAGE_ACCESS_TOKEN");
      localStorage.removeItem("FB_PAGE_ID");
    }
    setLoading(false);
  };

  const deleteConnection = () => {
    setLoading(true);
    localStorage.removeItem("FB_PAGE_DETAILS");
    localStorage.removeItem("FB_ACCESS_TOKEN");
    localStorage.removeItem("FB_PAGE_ACCESS_TOKEN");
    localStorage.removeItem("FB_PAGE_ID");
    setIsConnected(false);
    setLoading(false);
  };

  const checkPageConnected = () => {
    const accessToken = localStorage.getItem("FB_ACCESS_TOKEN");
    if (
      accessToken !== null &&
      accessToken !== undefined &&
      accessToken !== ""
    ) {
      setIsConnected(true);
    } else {
      setIsConnected(false);
    }
  };

  useEffect(() => {
    document.title = "Connect Page - Richpanel Assessment";
    checkPageConnected();
  }, []);

  // useEffect(() => {
  //   console.log("Started use effect");
  //   initFacebookSdk().then(() => {
  //     getFacebookLoginStatus().then((response) => {
  //       if (response == null) {
  //         console.log("No login status for the person");
  //       } else {
  //         console.log(response);
  //       }
  //     });
  //   });
  // }, []);

  const login = () => {
    console.log("reached log in button");
    fbLogin().then((response) => {
      console.log(response);
      if (response.status === "connected") {
        console.log("Person is connected");
      } else {
        // something
      }
    });
  }

  console.log(facebookAppID, facebookRedirectURI)

  return (
    <div className="h-[100vh] w-[100vw] bg-primary flex justify-center items-center">
      <Card>
        <div className="flex flex-col items-center justify-center w-[300px] gap-7">
          <h1 className="font-semibold text-lg">Facebook Page Integration</h1>
          {isConnected ? (
            <>
              <div className="w-full flex flex-col gap-3">
                <Button
                  onClick={deleteConnection}
                  variant="DANGER"
                  loading={loading}
                >
                  Delete Integration
                </Button>
                <Link className="w-full" to="/helpdesk">
                  <Button className="w-full">Reply To Messages</Button>
                </Link>
              </div>
            </>
          ) : (
            // <Button className="w-full">Connect Page</Button>
            <FacebookLogin
              appId={facebookAppID}
              redirectUri={facebookRedirectURI}
              scope="pages_show_list,pages_messaging,pages_read_engagement,pages_manage_metadata"
              callback={responseFacebook}
              onFailure={() => {
                showError("Could not connect to the Facebook page");
              }}
              render={(renderProps) => (
                <Button
                  onClick={renderProps.onClick}
                  loading={loading}
                  className="w-full"
                >
                  Connect Page
                </Button>
              )}
            />
            // <p>hello</p>
                // <Button
                //   onClick={login}
                //   loading={loading}
                //   className="w-full"
                // >
                //   Connect Page
                // </Button>
          )}
        </div>
      </Card>
    </div>
  );
};

export default FBConnectPage;
