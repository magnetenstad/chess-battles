from __future__ import annotations

import random
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional

import pygame

try:
    import chess
except ModuleNotFoundError as exc:
    raise SystemExit(
        "Missing dependency: python-chess. Install it with `pip install python-chess`."
    ) from exc

WINDOW_WIDTH = 1260
WINDOW_HEIGHT = 900
SQUARE_SIZE = 72
BOARD_SIZE = 8 * SQUARE_SIZE
LEFT_MARGIN = 30
TOP_MARGIN = 110
BOARD_GAP = 48
FPS = 60
BOT_DELAY_RANGE_MS = (450, 950)
NOTIFICATION_DURATION_MS = 3800
ATLAS_TILE_SIZE = 16
PIECE_RENDER_SIZE = 64

WHITE = chess.WHITE
BLACK = chess.BLACK

BG_COLOR = (22, 26, 34)
PANEL_COLOR = (32, 38, 50)
LIGHT_SQUARE = (236, 218, 185)
DARK_SQUARE = (181, 136, 99)
SELECTED_SQUARE = (95, 170, 120)
TARGET_SQUARE = (100, 155, 210)
BORDER_COLOR = (15, 18, 24)
TEXT_COLOR = (240, 240, 240)
MUTED_TEXT = (194, 202, 219)
ALERT_TEXT = (248, 129, 129)
SUCCESS_TEXT = (145, 232, 168)

PIECE_VALUES = {
    chess.PAWN: 1,
    chess.KNIGHT: 3,
    chess.BISHOP: 3,
    chess.ROOK: 5,
    chess.QUEEN: 9,
    chess.KING: 0,
}


@dataclass
class BoardState:
    index: int
    origin: tuple[int, int]
    label: str
    board: chess.Board = field(default_factory=chess.Board)
    selected_square: Optional[int] = None
    legal_targets: set[int] = field(default_factory=set)
    next_bot_move_ms: int = 0


class PawnSpawnPrototype:
    def __init__(self) -> None:
        pygame.init()
        pygame.display.set_caption("Dual Chess: Cross-Board Pawn Spawning")
        self.screen = pygame.display.set_mode((WINDOW_WIDTH, WINDOW_HEIGHT))
        self.clock = pygame.time.Clock()

        self.title_font = pygame.font.SysFont("georgia", 38, bold=True)
        self.header_font = pygame.font.SysFont("georgia", 25, bold=True)
        self.text_font = pygame.font.SysFont("consolas", 24, bold=True)
        self.small_font = pygame.font.SysFont("consolas", 19)
        self.piece_sprites = self.load_piece_sprites()

        self.rng = random.Random()
        self.running = True
        self.reset_game()

    def reset_game(self) -> None:
        left_origin = (LEFT_MARGIN, TOP_MARGIN)
        right_origin = (LEFT_MARGIN + BOARD_SIZE + BOARD_GAP, TOP_MARGIN)

        self.boards = [
            BoardState(index=0, origin=left_origin, label="Board 1: White A vs Bot A"),
            BoardState(index=1, origin=right_origin, label="Board 2: White B vs Bot B"),
        ]

        self.notifications: list[tuple[str, int]] = []
        self.game_over_message: Optional[str] = None

    def load_piece_sprites(self) -> dict[tuple[bool, int], pygame.Surface]:
        atlas_path = (
            Path(__file__).resolve().parent.parent
            / "assets"
            / "roupiks"
            / "atlas.png"
        )
        try:
            atlas = pygame.image.load(str(atlas_path)).convert_alpha()
        except pygame.error as exc:
            raise SystemExit(
                f"Could not load chess piece sprites from {atlas_path}"
            ) from exc

        sprite_cells = {
            (BLACK, chess.KING): (1, 13),
            (BLACK, chess.QUEEN): (2, 13),
            (BLACK, chess.ROOK): (3, 13),
            (BLACK, chess.BISHOP): (4, 13),
            (BLACK, chess.KNIGHT): (5, 13),
            (BLACK, chess.PAWN): (6, 13),
            (WHITE, chess.KING): (1, 14),
            (WHITE, chess.QUEEN): (2, 14),
            (WHITE, chess.ROOK): (3, 14),
            (WHITE, chess.BISHOP): (4, 14),
            (WHITE, chess.KNIGHT): (5, 14),
            (WHITE, chess.PAWN): (6, 14),
        }

        sprites: dict[tuple[bool, int], pygame.Surface] = {}
        for piece_key, (tile_x, tile_y) in sprite_cells.items():
            src = pygame.Rect(
                tile_x * ATLAS_TILE_SIZE,
                tile_y * ATLAS_TILE_SIZE,
                ATLAS_TILE_SIZE,
                ATLAS_TILE_SIZE,
            )
            piece_sprite = atlas.subsurface(src).copy()
            sprites[piece_key] = pygame.transform.scale(
                piece_sprite, (PIECE_RENDER_SIZE, PIECE_RENDER_SIZE)
            )
        return sprites

    def add_notification(self, text: str, duration_ms: int = NOTIFICATION_DURATION_MS) -> None:
        expires = pygame.time.get_ticks() + duration_ms
        self.notifications.append((text, expires))
        self.notifications = self.notifications[-8:]

    def clear_selection(self, state: BoardState) -> None:
        state.selected_square = None
        state.legal_targets.clear()

    def square_from_mouse(self, mouse_pos: tuple[int, int]) -> tuple[Optional[BoardState], Optional[int]]:
        mx, my = mouse_pos
        for state in self.boards:
            ox, oy = state.origin
            if ox <= mx < ox + BOARD_SIZE and oy <= my < oy + BOARD_SIZE:
                file_idx = (mx - ox) // SQUARE_SIZE
                rank_from_top = (my - oy) // SQUARE_SIZE
                board_rank = 7 - rank_from_top
                square = chess.square(int(file_idx), int(board_rank))
                return state, square
        return None, None

    def select_square(self, state: BoardState, square: int) -> None:
        piece = state.board.piece_at(square)
        if piece is None or piece.color != WHITE:
            self.clear_selection(state)
            return

        state.selected_square = square
        state.legal_targets = {
            move.to_square for move in state.board.legal_moves if move.from_square == square
        }

    def choose_white_move(self, candidates: list[chess.Move]) -> chess.Move:
        for move in candidates:
            if move.promotion == chess.QUEEN:
                return move
        for move in candidates:
            if move.promotion is None:
                return move
        return candidates[0]

    def handle_click(self, mouse_pos: tuple[int, int]) -> None:
        if self.game_over_message:
            return

        state, square = self.square_from_mouse(mouse_pos)
        if state is None or square is None:
            return

        board = state.board
        if board.is_game_over():
            return

        if board.turn != WHITE:
            self.add_notification(f"Board {state.index + 1}: waiting for black bot...")
            return

        clicked_piece = board.piece_at(square)

        if state.selected_square is None:
            if clicked_piece and clicked_piece.color == WHITE:
                self.select_square(state, square)
            return

        if square == state.selected_square:
            self.clear_selection(state)
            return

        candidates = [
            move
            for move in board.legal_moves
            if move.from_square == state.selected_square and move.to_square == square
        ]

        if candidates:
            chosen = self.choose_white_move(candidates)
            self.play_white_move(state, chosen)
            return

        if clicked_piece and clicked_piece.color == WHITE:
            self.select_square(state, square)
        else:
            self.clear_selection(state)

    def play_white_move(self, state: BoardState, move: chess.Move) -> None:
        board = state.board
        is_capture = board.is_capture(move)
        san = board.san(move)

        board.push(move)
        self.clear_selection(state)
        self.add_notification(f"Board {state.index + 1}: White played {san}")

        if is_capture:
            other = self.boards[1 - state.index]
            spawned = self.spawn_black_pawn(other)
            if spawned is None:
                self.add_notification(
                    f"Capture on Board {state.index + 1}: no empty square to spawn a pawn on Board {other.index + 1}"
                )
            else:
                square_name = chess.square_name(spawned)
                self.add_notification(
                    f"Capture on Board {state.index + 1}: spawned black pawn on Board {other.index + 1} at {square_name}"
                )

        self.check_for_game_over()
        if not self.game_over_message:
            self.ensure_bot_timers()

    def spawn_black_pawn(self, target_state: BoardState) -> Optional[int]:
        board = target_state.board
        preferred_ranks = [6, 5, 4, 3, 2, 1]

        chosen_square: Optional[int] = None
        for rank in preferred_ranks:
            rank_candidates: list[int] = []
            for file_idx in range(8):
                square = chess.square(file_idx, rank)
                if board.piece_at(square) is None:
                    rank_candidates.append(square)
            if rank_candidates:
                chosen_square = self.rng.choice(rank_candidates)
                break

        if chosen_square is None:
            return None

        board.set_piece_at(chosen_square, chess.Piece(chess.PAWN, BLACK))
        board.halfmove_clock = 0
        board.clear_stack()
        return chosen_square

    def evaluate_material(self, board: chess.Board) -> float:
        score = 0.0
        for piece in board.piece_map().values():
            value = PIECE_VALUES[piece.piece_type]
            score += value if piece.color == BLACK else -value
        return score

    def score_bot_move(self, board: chess.Board, move: chess.Move) -> float:
        score = 0.0

        attacker = board.piece_at(move.from_square)
        captured = board.piece_at(move.to_square)
        if board.is_capture(move):
            captured_value = PIECE_VALUES[chess.PAWN] if board.is_en_passant(move) else PIECE_VALUES[captured.piece_type] if captured else 0
            attacker_value = PIECE_VALUES[attacker.piece_type] if attacker else 0
            score += 12.0 * captured_value - 2.0 * attacker_value

        if board.gives_check(move):
            score += 2.5

        if move.promotion is not None:
            score += 8.0

        board.push(move)

        if board.is_checkmate():
            score += 100000.0
        else:
            score += self.evaluate_material(board) * 1.8
            if board.is_check():
                score += 3.0
            if board.is_stalemate():
                score -= 40.0

            moved_piece = board.piece_at(move.to_square)
            if moved_piece and board.is_attacked_by(WHITE, move.to_square):
                score -= PIECE_VALUES[moved_piece.piece_type] * 0.5

        board.pop()
        return score + self.rng.uniform(-0.2, 0.2)

    def choose_bot_move(self, board: chess.Board) -> Optional[chess.Move]:
        legal_moves = list(board.legal_moves)
        if not legal_moves:
            return None

        best_score = -sys.maxsize
        best_moves: list[chess.Move] = []

        for move in legal_moves:
            score = self.score_bot_move(board, move)
            if score > best_score + 1e-9:
                best_score = score
                best_moves = [move]
            elif abs(score - best_score) <= 1e-9:
                best_moves.append(move)

        return self.rng.choice(best_moves)

    def ensure_bot_timers(self) -> None:
        now = pygame.time.get_ticks()
        for state in self.boards:
            if state.board.is_game_over() or state.board.turn != BLACK:
                state.next_bot_move_ms = 0
                continue

            if state.next_bot_move_ms <= now:
                delay = self.rng.randint(*BOT_DELAY_RANGE_MS)
                state.next_bot_move_ms = now + delay

    def process_bot_turns(self) -> None:
        now = pygame.time.get_ticks()

        for state in self.boards:
            board = state.board
            if board.is_game_over() or board.turn != BLACK:
                continue

            if state.next_bot_move_ms == 0:
                delay = self.rng.randint(*BOT_DELAY_RANGE_MS)
                state.next_bot_move_ms = now + delay
                continue

            if now < state.next_bot_move_ms:
                continue

            move = self.choose_bot_move(board)
            state.next_bot_move_ms = 0
            if move is None:
                continue

            san = board.san(move)
            board.push(move)
            self.add_notification(f"Board {state.index + 1}: Black played {san}")

            self.check_for_game_over()
            if self.game_over_message:
                return

    def check_for_game_over(self) -> None:
        if self.game_over_message:
            return

        for state in self.boards:
            board = state.board
            if board.is_checkmate():
                if board.turn == BLACK:
                    self.game_over_message = (
                        f"Board {state.index + 1}: White checkmated Black. White team wins."
                    )
                else:
                    self.game_over_message = (
                        f"Board {state.index + 1}: Black checkmated White. White team loses."
                    )
                self.add_notification(self.game_over_message, duration_ms=8000)
                return

        if all(board_state.board.is_game_over() for board_state in self.boards):
            self.game_over_message = "Both boards ended without checkmate."
            self.add_notification(self.game_over_message, duration_ms=8000)

    def board_status_text(self, state: BoardState) -> tuple[str, tuple[int, int, int]]:
        board = state.board

        if board.is_checkmate():
            if board.turn == BLACK:
                return "Status: Black is checkmated", SUCCESS_TEXT
            return "Status: White is checkmated", ALERT_TEXT

        if board.is_stalemate():
            return "Status: Stalemate", MUTED_TEXT

        if board.is_insufficient_material():
            return "Status: Draw (insufficient material)", MUTED_TEXT

        turn_text = "White (human)" if board.turn == WHITE else "Black (bot)"
        if board.is_check():
            return f"Turn: {turn_text}  |  CHECK", ALERT_TEXT
        return f"Turn: {turn_text}", TEXT_COLOR

    def draw_piece(self, piece: chess.Piece, center: tuple[int, int]) -> None:
        sprite = self.piece_sprites[(piece.color, piece.piece_type)]
        shadow = pygame.Surface((PIECE_RENDER_SIZE, PIECE_RENDER_SIZE), pygame.SRCALPHA)
        shadow_rect = pygame.Rect(
            (PIECE_RENDER_SIZE - 12) // 2,
            PIECE_RENDER_SIZE - 18,
            PIECE_RENDER_SIZE - 12,
            10,
        )
        pygame.draw.ellipse(shadow, (0, 0, 0, 85), shadow_rect)
        shadow_blit = shadow.get_rect(center=(center[0], center[1] + PIECE_RENDER_SIZE // 3))
        self.screen.blit(shadow, shadow_blit)
        piece_rect = sprite.get_rect(center=center)
        self.screen.blit(sprite, piece_rect)

    def draw_board(self, state: BoardState) -> None:
        ox, oy = state.origin
        board = state.board

        for rank_from_top in range(8):
            board_rank = 7 - rank_from_top
            for file_idx in range(8):
                square = chess.square(file_idx, board_rank)
                rect = pygame.Rect(
                    ox + file_idx * SQUARE_SIZE,
                    oy + rank_from_top * SQUARE_SIZE,
                    SQUARE_SIZE,
                    SQUARE_SIZE,
                )

                is_dark = (file_idx + board_rank) % 2 == 0
                color = DARK_SQUARE if is_dark else LIGHT_SQUARE

                if square == state.selected_square:
                    color = SELECTED_SQUARE
                elif square in state.legal_targets:
                    color = TARGET_SQUARE

                pygame.draw.rect(self.screen, color, rect)

                if square in state.legal_targets and square != state.selected_square:
                    pygame.draw.circle(
                        self.screen,
                        (40, 55, 75),
                        rect.center,
                        SQUARE_SIZE // 8,
                    )

                piece = board.piece_at(square)
                if piece:
                    self.draw_piece(piece, rect.center)

        border_rect = pygame.Rect(ox, oy, BOARD_SIZE, BOARD_SIZE)
        pygame.draw.rect(self.screen, BORDER_COLOR, border_rect, width=4)

        for idx, file_name in enumerate("abcdefgh"):
            label = self.small_font.render(file_name, True, MUTED_TEXT)
            lx = ox + idx * SQUARE_SIZE + SQUARE_SIZE // 2 - label.get_width() // 2
            self.screen.blit(label, (lx, oy + BOARD_SIZE + 4))

        for idx, rank_name in enumerate("87654321"):
            label = self.small_font.render(rank_name, True, MUTED_TEXT)
            ly = oy + idx * SQUARE_SIZE + SQUARE_SIZE // 2 - label.get_height() // 2
            self.screen.blit(label, (ox - 18, ly))

        header = self.header_font.render(state.label, True, TEXT_COLOR)
        self.screen.blit(header, (ox, oy - 42))

        status_text, status_color = self.board_status_text(state)
        status = self.small_font.render(status_text, True, status_color)
        self.screen.blit(status, (ox, oy + BOARD_SIZE + 28))

    def draw_notifications(self) -> None:
        now = pygame.time.get_ticks()
        self.notifications = [(text, exp) for text, exp in self.notifications if exp > now]

        panel = pygame.Rect(LEFT_MARGIN, WINDOW_HEIGHT - 170, WINDOW_WIDTH - 2 * LEFT_MARGIN, 136)
        pygame.draw.rect(self.screen, PANEL_COLOR, panel, border_radius=10)
        pygame.draw.rect(self.screen, BORDER_COLOR, panel, width=2, border_radius=10)

        title = self.small_font.render("Event Log", True, MUTED_TEXT)
        self.screen.blit(title, (panel.x + 12, panel.y + 8))

        start_y = panel.y + 34
        for row, (text, _) in enumerate(self.notifications[-4:]):
            line = self.small_font.render(text, True, TEXT_COLOR)
            self.screen.blit(line, (panel.x + 12, start_y + row * 24))

    def draw(self) -> None:
        self.screen.fill(BG_COLOR)

        title = self.title_font.render("Dual Chess: Captures Spawn Enemy Pawns", True, TEXT_COLOR)
        self.screen.blit(title, (LEFT_MARGIN, 22))

        instruction = self.small_font.render(
            "Click to move White on either board. Each White capture spawns a Black pawn on the other board.",
            True,
            MUTED_TEXT,
        )
        self.screen.blit(instruction, (LEFT_MARGIN, 72))

        for state in self.boards:
            self.draw_board(state)

        self.draw_notifications()

        if self.game_over_message:
            overlay = pygame.Surface((WINDOW_WIDTH, WINDOW_HEIGHT), pygame.SRCALPHA)
            overlay.fill((0, 0, 0, 145))
            self.screen.blit(overlay, (0, 0))

            msg = self.header_font.render(self.game_over_message, True, SUCCESS_TEXT)
            msg_rect = msg.get_rect(center=(WINDOW_WIDTH // 2, WINDOW_HEIGHT // 2 - 20))
            self.screen.blit(msg, msg_rect)

            tip = self.small_font.render("Press R to restart or Esc to quit.", True, TEXT_COLOR)
            tip_rect = tip.get_rect(center=(WINDOW_WIDTH // 2, WINDOW_HEIGHT // 2 + 18))
            self.screen.blit(tip, tip_rect)

        pygame.display.flip()

    def run(self) -> None:
        while self.running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    self.running = False
                elif event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_ESCAPE:
                        self.running = False
                    elif event.key == pygame.K_r:
                        self.reset_game()
                elif event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
                    self.handle_click(event.pos)

            if not self.game_over_message:
                self.process_bot_turns()

            self.draw()
            self.clock.tick(FPS)

        pygame.quit()


if __name__ == "__main__":
    PawnSpawnPrototype().run()
