const API_BASE = "http://localhost:8080";
const SQUARE_SIZE = 70;
const LIGHT = "#3a352c";
const DARK = "#201e19";
const HIGHLIGHT = "rgba(217, 154, 58, 0.45)";
const SELECTED = "rgba(217, 154, 58, 0.75)";

const PIECE_GLYPHS = {
  P: "♙", N: "♘", B: "♗", R: "♖", Q: "♕", K: "♔",
  p: "♟", n: "♞", b: "♝", r: "♜", q: "♛", k: "♚",
};

const canvas = document.getElementById("board");
const ctx = canvas.getContext("2d");
const statusEl = document.getElementById("status");
const moveListEl = document.getElementById("moveList");
const newGameBtn = document.getElementById("newGameBtn");
const difficultySelect = document.getElementById("difficulty");
const promoModal = document.getElementById("promoModal");
const capturedByBlackEl = document.getElementById("capturedByBlack");
const capturedByWhiteEl = document.getElementById("capturedByWhite");

let gameId = null;
let boardState = null;
let selectedSquare = null;
let legalTargets = [];
let pendingPromotion = null;

function fenToBoard(fen) {
  const rows = fen.split(" ")[0].split("/");
  const grid = [];
  for (const row of rows) {
    const line = [];
    for (const ch of row) {
      if (/\d/.test(ch)) {
        for (let i = 0; i < Number(ch); i++) line.push(null);
      } else {
        line.push(ch);
      }
    }
    grid.push(line);
  }
  return grid;
}

function squareToCoords(file, rank) {
  const x = file * SQUARE_SIZE;
  const y = (7 - rank) * SQUARE_SIZE;
  return { x, y };
}

function algebraicToFileRank(sq) {
  const file = sq.charCodeAt(0) - "a".charCodeAt(0);
  const rank = Number(sq[1]) - 1;
  return { file, rank };
}

function fileRankToAlgebraic(file, rank) {
  return String.fromCharCode("a".charCodeAt(0) + file) + (rank + 1);
}

function drawBoard() {
  if (!boardState) return;
  const grid = fenToBoard(boardState.fen);

  for (let rank = 0; rank < 8; rank++) {
    for (let file = 0; file < 8; file++) {
      const { x, y } = squareToCoords(file, rank);
      const isLight = (file + rank) % 2 === 1;
      ctx.fillStyle = isLight ? LIGHT : DARK;
      ctx.fillRect(x, y, SQUARE_SIZE, SQUARE_SIZE);

      const alg = fileRankToAlgebraic(file, rank);
      if (selectedSquare === alg) {
        ctx.fillStyle = SELECTED;
        ctx.fillRect(x, y, SQUARE_SIZE, SQUARE_SIZE);
      } else if (legalTargets.includes(alg)) {
        ctx.fillStyle = HIGHLIGHT;
        ctx.beginPath();
        ctx.arc(x + SQUARE_SIZE / 2, y + SQUARE_SIZE / 2, 12, 0, Math.PI * 2);
        ctx.fill();
      }

      const piece = grid[7 - rank][file];
      if (piece) {
        ctx.fillStyle = piece === piece.toUpperCase() ? "#f0ead9" : "#1a1712";
        ctx.strokeStyle = piece === piece.toUpperCase() ? "#1a1712" : "#f0ead9";
        ctx.lineWidth = 1;
        ctx.font = "48px serif";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        const glyph = PIECE_GLYPHS[piece];
        ctx.fillText(glyph, x + SQUARE_SIZE / 2, y + SQUARE_SIZE / 2 + 2);
      }
    }
  }
}

function squareFromEvent(evt) {
  const rect = canvas.getBoundingClientRect();
  const x = evt.clientX - rect.left;
  const y = evt.clientY - rect.top;
  const file = Math.floor(x / SQUARE_SIZE);
  const rank = 7 - Math.floor(y / SQUARE_SIZE);
  if (file < 0 || file > 7 || rank < 0 || rank > 7) return null;
  return fileRankToAlgebraic(file, rank);
}

function legalTargetsFrom(square) {
  return boardState.legal_moves
    .filter((m) => m.startsWith(square))
    .map((m) => m.slice(2, 4));
}

async function newGame() {
  const res = await fetch(`${API_BASE}/api/new-game`, { method: "POST" });
  const data = await res.json();
  gameId = data.game_id;
  selectedSquare = null;
  legalTargets = [];
  moveListEl.innerHTML = "";
  await refreshState();
}

async function refreshState() {
  const res = await fetch(`${API_BASE}/api/state?game_id=${gameId}`);
  boardState = await res.json();
  drawBoard();
  updateStatus();
  updateCaptures();
}

function updateStatus() {
  if (!boardState) return;
  if (boardState.game_over) {
    statusEl.textContent = boardState.in_check
      ? `Checkmate — ${boardState.side_to_move} loses`
      : "Stalemate";
    return;
  }
  statusEl.textContent = boardState.in_check
    ? `${capitalize(boardState.side_to_move)} to move — check!`
    : `${capitalize(boardState.side_to_move)} to move`;
}

function updateCaptures() {
  capturedByBlackEl.textContent = (boardState.captured_by_black || [])
    .map((p) => PIECE_GLYPHS[p])
    .join(" ");
  capturedByWhiteEl.textContent = (boardState.captured_by_white || [])
    .map((p) => PIECE_GLYPHS[p])
    .join(" ");
}

function capitalize(s) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function logMove(moveStr) {
  const li = document.createElement("li");
  li.textContent = moveStr;
  moveListEl.appendChild(li);
  moveListEl.scrollTop = moveListEl.scrollHeight;
}

async function playerMove(from, to, promotion) {
  const res = await fetch(`${API_BASE}/api/move`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ game_id: gameId, from, to, promotion: promotion || "" }),
  });
  if (!res.ok) {
    selectedSquare = null;
    legalTargets = [];
    drawBoard();
    return;
  }
  const data = await res.json();
  boardState = data.state;
  logMove(data.move);
  selectedSquare = null;
  legalTargets = [];
  drawBoard();
  updateStatus();
  updateCaptures();

  if (!boardState.game_over) {
    await engineMove();
  }
}

async function engineMove() {
  statusEl.textContent = "Engine thinking...";
  const thinkMs = Number(difficultySelect.value);
  const res = await fetch(`${API_BASE}/api/engine-move`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ game_id: gameId, think_ms: thinkMs }),
  });
  const data = await res.json();
  boardState = data.state;
  logMove(data.move);
  drawBoard();
  updateStatus();
  updateCaptures();
}

function isPromotionMove(from, to) {
  const toRank = Number(to[1]);
  const fromRank = Number(from[1]);
  const piece = pieceAt(from);
  return piece && piece.toLowerCase() === "p" && (toRank === 8 || toRank === 1) && fromRank !== toRank;
}

function pieceAt(square) {
  const grid = fenToBoard(boardState.fen);
  const { file, rank } = algebraicToFileRank(square);
  return grid[7 - rank][file];
}

canvas.addEventListener("click", (evt) => {
  if (!boardState || boardState.game_over) return;
  const square = squareFromEvent(evt);
  if (!square) return;

  if (selectedSquare && legalTargets.includes(square)) {
    if (isPromotionMove(selectedSquare, square)) {
      pendingPromotion = { from: selectedSquare, to: square };
      promoModal.classList.remove("hidden");
      return;
    }
    playerMove(selectedSquare, square, "");
    return;
  }

  const piece = pieceAt(square);
  if (piece && piece === piece.toUpperCase()) {
    selectedSquare = square;
    legalTargets = legalTargetsFrom(square);
  } else {
    selectedSquare = null;
    legalTargets = [];
  }
  drawBoard();
});

promoModal.addEventListener("click", (evt) => {
  const promo = evt.target.getAttribute("data-promo");
  if (!promo || !pendingPromotion) return;
  promoModal.classList.add("hidden");
  playerMove(pendingPromotion.from, pendingPromotion.to, promo);
  pendingPromotion = null;
});

newGameBtn.addEventListener("click", newGame);

newGame();
