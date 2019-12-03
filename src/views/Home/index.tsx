import React from 'react';
import { withAuth } from '@okta/okta-react';
import useAuth from 'hooks/useAuth';

type HomeProps = {
  auth: any;
};

const Home = ({ auth }: HomeProps) => {
  const [isAuthenticated] = useAuth(auth);
  return (
    <div>
      <h1>Home</h1>
      <h3>{`A user is ${isAuthenticated ? '' : 'NOT'} authenticated`}</h3>
    </div>
  );
};

export default withAuth(Home);
