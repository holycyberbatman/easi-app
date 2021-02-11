import React, { useEffect, useState } from 'react';

import Footer from 'components/Footer';
import Header from 'components/Header';
import MainContent from 'components/MainContent';
import PageWrapper from 'components/PageWrapper';
import OktaSignInWidget from 'components/shared/OktaSignInWidget';
import { localAuthStorageKey } from 'constants/localAuth';
import usePageContext from 'hooks/usePageContext';
import { isLocalEnvironment } from 'utils/local';
import DevLogin from 'views/AuthenticationWrapper/DevLogin';

const Login = () => {
  const { setPage } = usePageContext();

  let defaultAuth = false;
  if (isLocalEnvironment() && window.localStorage[localAuthStorageKey]) {
    defaultAuth = JSON.parse(window.localStorage[localAuthStorageKey])
      .favorLocalAuth;
  }
  const [isLocalAuth, setIsLocalAuth] = useState(defaultAuth);

  const handleUseLocalAuth = () => {
    setIsLocalAuth(true);
  };

  useEffect(() => {
    setPage('login page');
  }, [setPage]);

  if (isLocalEnvironment() && isLocalAuth) {
    return (
      <PageWrapper>
        <MainContent className="grid-container margin-top-4">
          <DevLogin />
        </MainContent>
      </PageWrapper>
    );
  }

  return (
    <PageWrapper>
      <Header />
      <MainContent className="grid-container">
        {isLocalEnvironment() && (
          <div>
            <button type="button" onClick={handleUseLocalAuth}>
              Use Local Auth
            </button>
          </div>
        )}
        <OktaSignInWidget onSuccess={() => {}} onError={() => {}} />
      </MainContent>
      <Footer />
    </PageWrapper>
  );
};

export default Login;
