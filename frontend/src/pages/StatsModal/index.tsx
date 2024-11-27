import {
  Box,
  Button,
  Typography,
  Modal,
  Divider,
} from "@mui/material";


type Stat = {
  total_numberof_requests: number,
  succeeded_requests: number,
  failed_requests: number,
  target_users: number,
}
interface Props {
  open: boolean;
  stats: Stat
  setOpen: (val:boolean) => void
}

export default function StatsModal({open,stats,setOpen}:Props) {


  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  return (
    <>
      <Button variant="contained" color="primary" onClick={handleOpen}>
        View Stats
      </Button>

      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="modal-title"
        aria-describedby="modal-description"
      >
        <Box
          sx={{
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            width: 400,
            bgcolor: "background.paper",
            boxShadow: 24,
            borderRadius: 2,
            p: 4,
          }}
        >
          <Typography
            id="modal-title"
            variant="h6"
            component="h2"
            sx={{ mb: 2, textAlign: "center" }}
          >
            Statistics
          </Typography>
          <Divider sx={{ mb: 2 }} />
          <Typography id="modal-description" variant="body1" sx={{ mb: 1 }}>
            <strong>Total Number of Requests:</strong> {stats.total_numberof_requests}
          </Typography>
          <Typography variant="body1" sx={{ mb: 1 }}>
            <strong>Succeeded Requests:</strong> {stats.succeeded_requests}
          </Typography>
          <Typography variant="body1" sx={{ mb: 1 }}>
            <strong>Failed Requests:</strong> {stats.failed_requests}
          </Typography>
          <Typography variant="body1" sx={{ mb: 1 }}>
            <strong>Target Users:</strong> {stats.target_users}
          </Typography>
          <Divider sx={{ mt: 2, mb: 2 }} />
          <Box sx={{ textAlign: "center" }}>
            <Button variant="contained" color="secondary" onClick={handleClose}>
              Close
            </Button>
          </Box>
        </Box>
      </Modal>
    </>
  );
}
