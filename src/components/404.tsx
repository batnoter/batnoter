import React from "react";
import ErrorImage from "../assets/404.png";
import { Grid, Typography, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";

const ErrorPage: React.FC = (): React.ReactElement => {
  const history = useNavigate();
  const handleClick = () => history("/");
  return (
    <Grid container columns={12} minHeight="100vh">
      <Grid
        item
        xs={12}
        md={5}
        lg={5}
        sx={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Typography variant="h1" sx={{ fontWeight: "bold" }}>
          Whooops!
        </Typography>
        <Typography variant="body1">
          Sorry, the Page you are looking for does not exist.
        </Typography>
        <Button
          variant="contained"
          sx={{ mt: 2, width: "20em", alignSelf: "left" }}
          onClick={handleClick}
        >
          Go back to Home
        </Button>
      </Grid>
      <Grid
        item
        xs={12}
        md={7}
        lg={7}
        sx={{ display: "grid", placeContent: "center" }}
      >
        <img
          src={ErrorImage}
          alt="404"
          width="100%"
          style={{ objectFit: "cover" }}
        />
      </Grid>
    </Grid>
  );
};

export default ErrorPage;
