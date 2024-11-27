import React, { useEffect, useState } from "react";
import { Box, CircularProgress, Typography, Grid, Paper, Divider } from "@mui/material";
import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  Legend,
} from "recharts";
import { useParams } from "react-router-dom";
import toast from "react-hot-toast";

export interface Test {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  uuid: string;
  url: string;
  method: string;
  headers: object;
  target_users: number;
  reach_peak_after_in_minutes: number;
  users_to_start_with: number;
  status: string;
  total_requests: number;
  succeeded_requests: number;
  report: Report;
}

export interface Report {
  average_response_time: number;
  peak_response_time: number;
  error_rate: number;
  throughput: number;
  p_50_percentile: number;
  p_90_percentile: number;
  p_99_percentile: number;
}

export default function MetricsPage() {
  const { id } = useParams();
  const [state, setState] = useState<Test | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const API_URL = import.meta.env.VITE_API_URL;

  useEffect(() => {
    if (id) {
      setIsLoading(true);
      fetch(`${API_URL}/tests/${id}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((res) => res.json())
        .then((data) => {
          setState(data);
        })
        .catch((err) => {
          console.error("Error in fetching test report", err);
          toast.error("Failed to load metrics");
        })
        .finally(() => {
          setIsLoading(false);
        });
    }
  }, [id]);

  const performanceData = state?.report
    ? [
        { metric: "Average Response Time", value: state.report.average_response_time },
        { metric: "Peak Response Time", value: state.report.peak_response_time },
        { metric: "Error Rate", value: state.report.error_rate },
        { metric: "Throughput", value: state.report.throughput },
        { metric: "50th Percentile", value: state.report.p_50_percentile },
        { metric: "90th Percentile", value: state.report.p_90_percentile },
        { metric: "99th Percentile", value: state.report.p_99_percentile },
      ]
    : [];

  const requestsData = state
    ? [
        { type: "Total Requests", value: state.total_requests },
        { type: "Succeeded Requests", value: state.succeeded_requests },
        { type: "Failed Requests", value: state.total_requests - state.succeeded_requests },
      ]
    : [];

  return (
    <Box sx={{ padding: 4 }}>
      <Typography variant="h4" gutterBottom>
        Test Metrics
      </Typography>
      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", height: 300 }}>
          <CircularProgress />
        </Box>
      ) : state ? (
        <Grid container spacing={4}>
          {/* Key Details Section */}
          <Grid item xs={12} md={4}>
            <Paper sx={{ padding: 2 }}>
              <Typography variant="h6" gutterBottom>
                Test Details
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Typography variant="body1">
                <strong>URL:</strong> {state.url}
              </Typography>
              <Typography variant="body1">
                <strong>Method:</strong> {state.method}
              </Typography>
              <Typography variant="body1">
                <strong>Status:</strong> {state.status}
              </Typography>
              <Typography variant="body1">
                <strong>Total Requests:</strong> {state.total_requests}
              </Typography>
              <Typography variant="body1">
                <strong>Succeeded Requests:</strong> {state.succeeded_requests}
              </Typography>
              <Typography variant="body1">
                <strong>Failed Requests:</strong>{" "}
                {state.total_requests - state.succeeded_requests}
              </Typography>
            </Paper>
          </Grid>

          {/* Performance Metrics Graph */}
          <Grid item xs={12} md={8}>
            <Paper sx={{ padding: 2 }}>
              <Typography variant="h6" gutterBottom>
                Performance Metrics
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={performanceData} margin={{ top: 20, right: 30, left: 0, bottom: 5 }}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="metric" />
                  <YAxis />
                  <Tooltip />
                  <Line type="monotone" dataKey="value" stroke="#1976d2" strokeWidth={2} />
                </LineChart>
              </ResponsiveContainer>
            </Paper>
          </Grid>

          {/* Requests Graph */}
          <Grid item xs={12}>
            <Paper sx={{ padding: 2 }}>
              <Typography variant="h6" gutterBottom>
                Requests Breakdown
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={requestsData} margin={{ top: 20, right: 30, left: 0, bottom: 5 }}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="type" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="value" fill="#1976d2" />
                </BarChart>
              </ResponsiveContainer>
            </Paper>
          </Grid>
        </Grid>
      ) : (
        <Typography variant="body1" color="error">
          No data available.
        </Typography>
      )}
    </Box>
  );
}
