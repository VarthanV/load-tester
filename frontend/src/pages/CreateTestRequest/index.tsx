/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useState } from "react";
import { Box, TextField, Button, MenuItem, Typography, CircularProgress } from "@mui/material";

export default function CreateTestRequest() {
  const [formData, setFormData] = useState({
    url: "https://google.com",
    method: "GET",
    target_users: 100,
    users_to_start_with: 1,
    reach_peak_afer_in_minutes: 2,
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const API_URL = import.meta.env.VITE_API_URL; 

  const handleChange = (e: any) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await fetch(`${API_URL}/test`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} - ${response.statusText}`);
      }

      const data = await response.json();
      setSuccess("Request submitted successfully!");
      console.log("Response Data:", data);
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Box
        component="form"
        sx={{
          maxWidth: 600,
          margin: "auto",
          padding: 3,
          display: "flex",
          flexDirection: "column",
          gap: 2,
          boxShadow: 3,
          borderRadius: 2,
          marginTop: "3rem",
        }}
        onSubmit={handleSubmit}
      >
        <Typography variant="h5" align="center">
          Test Configuration
        </Typography>

        <TextField
          label="URL"
          name="url"
          value={formData.url}
          onChange={handleChange}
          fullWidth
          required
        />

        <TextField
          select
          label="Method"
          name="method"
          value={formData.method}
          onChange={handleChange}
          fullWidth
          required
        >
          {["GET", "POST", "PUT", "DELETE"].map((method) => (
            <MenuItem key={method} value={method}>
              {method}
            </MenuItem>
          ))}
        </TextField>

        <TextField
          label="Target Users"
          name="target_users"
          type="number"
          value={formData.target_users}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ min: 1 }}
        />

        <TextField
          label="Users to Start With"
          name="users_to_start_with"
          type="number"
          value={formData.users_to_start_with}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ min: 1 }}
        />

        <TextField
          label="Reach Peak After (in minutes)"
          name="reach_peak_afer_in_minutes"
          type="number"
          value={formData.reach_peak_afer_in_minutes}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ min: 1 }}
        />

        {loading ? (
          <CircularProgress color="secondary" />
        ) : (
          <Button type="submit" variant="contained" color="primary" fullWidth>
            Submit
          </Button>
        )}

        {success && <Typography color="green">{success}</Typography>}
        {error && <Typography color="red">{error}</Typography>}
      </Box>
    </>
  );
}
