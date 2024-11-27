import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import { useNavigate } from "react-router-dom";

export default function Navbar() {
  const navigate = useNavigate();
  return (
    <Box component={"section"} sx={{ flexGrow: 1, width: "100%" }}>
      <AppBar position="static" sx={{ width: "100%" }}>
        <Toolbar>
          <Typography
            variant="h6"
            component="div"
            sx={{ flexGrow: 1 }}
            onClick={(e) => {
              e.preventDefault();
              navigate("/");
            }}
          >
            ❤️‍🔥 Load Tester
          </Typography>
          <Button
            color="inherit"
            onClick={(e) => {
              e.preventDefault();
              navigate("/test");
            }}
          >
            New
          </Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
