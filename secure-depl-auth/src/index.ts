import express, { Request, Response } from "express";
import bodyParser from "body-parser";

const app = express();
const PORT = process.env.PORT || 3000;
const AUTH_TIMEOUT = 7000; // 7 seconds

let authStatus = { authorizedBy: "" };

// Middleware
app.use(bodyParser.json());

// POST /auth
app.post("/auth", (req: Request, res: Response) => {
  const name = req.body.name;

  if (name) {
    authStatus = { authorizedBy: name };
    setTimeout(() => (authStatus = { authorizedBy: "" }), AUTH_TIMEOUT);
    return res.status(200).json({ message: "Authorized" });
  }

  return res.status(400).json({ error: "Bad Request: name is required" });
});

// GET /auth
app.get("/auth", (req: Request, res: Response) => {
  if (!authStatus.authorizedBy) {
    return res.status(401).json({ error: "Unnauthorized" });
  }

  return res.status(200).json({ message: "Authorized" });
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
