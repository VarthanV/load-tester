/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useState, useEffect } from "react";
import {
  Box,
  TextField,
  Button,
  MenuItem,
  Typography,
  CircularProgress,
} from "@mui/material";
import { useNavigate } from "react-router-dom";

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
  const [testId, setTestId] = useState<string | null>(null);
  const [update, setUpdate] = useState<any>(null);
  const [isPolling, setIsPolling] = useState(false);

  const API_URL = import.meta.env.VITE_API_URL;
  const navigate = useNavigate();

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
      const response = await fetch(`${API_URL}/tests`, {
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
      setTestId(data.id); // Store the ID from the response
      setSuccess("Request submitted successfully, Polling updates!");
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!testId) return;

    const interval = setInterval(async () => {
      try {
        setIsPolling(true);
        const response = await fetch(`${API_URL}/tests/${testId}/updates`, {
          method: "GET",
        });

        if (!response.ok) {
          throw new Error(`Error: ${response.status} - ${response.statusText}`);
        }

        const data = await response.json();
        setUpdate(data.update);

        if (data.update === null) {
          clearInterval(interval);
          navigate(`/test/${testId}`);
        }
      } catch (err) {
        console.error("Polling error:", err);
        clearInterval(interval);
      } finally {
        setIsPolling(false);
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [testId]);

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
        {testId && (
          <Typography variant="h5" className="mt">
            {" "}
            Real Time updates
          </Typography>
        )}

        {update ? (
          !isPolling ? (
            <Box>
              <Typography>
                Total Requests: {update.total_numberof_requests}
              </Typography>
              <Typography>Succeeded: {update.succeeded_requests}</Typography>
              <Typography>Failed: {update.failed_requests}</Typography>
              <Typography>Target Users: {update.target_users}</Typography>
            </Box>
          ) : (
            <CircularProgress />
          )
        ) : null}
      </Box>
    </>
  );
}
