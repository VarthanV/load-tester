import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom"; // React Router for navigation
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  CircularProgress,
} from "@mui/material";
import toast from "react-hot-toast";

export type Tests = Test[];

export interface Test {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  uuid: string;
  url: string;
  method: string;
  target_users: number;
  reach_peak_after_in_minutes: number;
  report: Report;
  users_to_start_with?: number;
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

const HomePage = () => {
  const API_URL = import.meta.env.VITE_API_URL;
  const navigate = useNavigate(); // Hook to navigate to a new route
  const [tests, setTests] = useState<Tests | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetch(`${API_URL}/tests`)
      .then((res) => res.json())
      .then((data) => {
        setTests(data);
      })
      .catch((err) => {
        toast.error("error in fetching tests");
        console.error("error in fetching tests ", err);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  const handleRowClick = (uuid: string) => {
    navigate(`/test/${uuid}`);
  };

  if (isLoading && !tests) {
    return <CircularProgress />;
  }

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>UUID</TableCell>
            <TableCell>URL</TableCell>
            <TableCell>Method</TableCell>
            <TableCell>Target Users</TableCell>
            <TableCell>Users to Start With</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tests &&
            tests.map((row) => (
              <TableRow
                key={row.ID}
                hover
                onClick={() => handleRowClick(row.uuid)} // Handle row click
              >
                <TableCell>{row.ID}</TableCell>
                <TableCell>{row.uuid}</TableCell>
                <TableCell>{row.url}</TableCell>
                <TableCell>{row.method}</TableCell>
                <TableCell>{row.target_users}</TableCell>
                <TableCell>{row.users_to_start_with || "-"}</TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default HomePage;
