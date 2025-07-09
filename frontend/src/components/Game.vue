<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const RECONNECT_INTERVAL = 1500

const router = useRouter()

const connection = ref<'connected' | 'disconnected' | 'reconnecting'>('disconnected')
const socket = ref<WebSocket | null>(null)
let reconnectTimeout: ReturnType<typeof setTimeout>

const connectWebSocket = () => {
  connection.value = 'reconnecting'
  const url = import.meta.env.VITE_WS_URL
  if (!url) throw new Error('VITE_WS_URL is not defined')

  const ws = new WebSocket(url)
  socket.value = ws

  ws.onopen = () => {
    connection.value = 'connected'

    ws.send(
      JSON.stringify({
        type: 'playerJoin',
      }),
    )
  }

  ws.onmessage = (event) => {
    const packet = JSON.parse(event.data)
    const room = gameStore.room

    switch (packet.type) {
      case 'errorResponse':
        router.push(`/error?message=${encodeURIComponent(packet.data.message)}`)
        break

      case 'playerJoinResponse':
        gameStore.setGame(packet.data.game)
        gameStore.setPlayers(packet.data.players)
        gameStore.setRoom(packet.data.room)
        break

      case 'playerJoinedEvent':
        if (packet.data.id !== gameStore.authUser!.id) {
          gameStore.setPlayers([...gameStore.players, packet.data])
          toast(`${packet.data.username} has joined the lobby.`)
        }
        break

      case 'playerChatEvent':
        gameStore.addChat(packet.data)
        break

      case 'gameStartEvent':
        gameStore.setGame({ ...gameStore.game!, state: 'ingame' })
        gameStore.setRoom(packet.data.room)
        gameStore.setStartingIn(packet.data.startUnix - currentTimeInSeconds())
        toast('When the countdown hits 0 the game board will be revealed.')
        break

      case 'gameEndEvent':
        gameStore.setGame({ ...gameStore.game!, state: 'lobby' })
        break

      case 'gameRevealBoard':
        gameStore.setRoom(packet.data.room)
        toast('When you have found a route submit a bid in the sidebar.')
        break

      case 'playerBidEvent':
        if (packet.data.endUnix) {
          gameStore.setStartingIn(packet.data.endUnix - currentTimeInSeconds())
          toast(
            `Player made the first bid. You have ${packet.data.endUnix - currentTimeInSeconds()}s.`,
          )
        }

        gameStore.setRoom({
          ...gameStore.room!,
          ingameState: 'countdown',
          currentBids: {
            ...Object.fromEntries(
              Object.entries(gameStore.room?.currentBids ?? {}).filter(
                ([id]) => parseInt(id) !== packet.data.playerId,
              ),
            ),
            [packet.data.playerId]: packet.data.bid,
          },
        })
        break

      case 'gameStartVerification':
        gameStore.setStartingIn(packet.data.endUnix - currentTimeInSeconds())
        gameStore.setRoom({
          ...gameStore.room!,
          currentVerifyingPlayerId: packet.data.playerId,
          ingameState: 'verify',
          movesTaken: 0,
        })
        toast('Verification has started')
        break

      case 'playerVerifyMove':
        if (gameStore.room?.currentVerifyingPlayerId !== gameStore.authUser?.id) {
          gameStore.setRoom({
            ...gameStore.room!,
            currentRockets: {
              ...gameStore.room!.currentRockets,
              [packet.data.rocket]: getCompletedMove(
                gameStore.room!.board,
                gameStore.room!.currentRockets,
                packet.data.rocket,
                packet.data.direction,
              ),
            },
            movesTaken: gameStore.room!.movesTaken! + 1,
          })
        }
        break

      case 'playerVerifyReset':
        if (gameStore.room?.currentVerifyingPlayerId !== gameStore.authUser?.id) {
          gameStore.setRoom({
            ...gameStore.room!,
            currentRockets: gameStore.room!.restorableRockets,
            movesTaken: 0,
          })
        }
        break

      case 'playerVerifyCompleted':
        gameStore.setRoom({
          ...gameStore.room!,
          wins: packet.data.wins,
          restorableRockets: { ...gameStore.room!.currentRockets },
          ingameState: 'winner',
        })
        gameStore.setStartingIn(packet.data.endUnix - currentTimeInSeconds())
        gameStore.setHighlightMovement(undefined)
        break

      case 'gameVerifyNext':
        gameStore.setRoom({
          ...gameStore.room!,
          currentBids: packet.data.newCurrentBids,
          currentRockets: { ...gameStore.room!.restorableRockets },
          movesTaken: 0,
          currentVerifyingPlayerId: packet.data.newVerifyingPlayerId,
        })
        gameStore.setStartingIn(packet.data.endUnix - currentTimeInSeconds())
        gameStore.setHighlightMovement(undefined)
        break

      case 'gameVerifyFailed':
        gameStore.setRoom({
          ...gameStore.room!,
          currentBids: {},
          currentRockets: { ...gameStore.room!.restorableRockets },
          currentVerifyingPlayerId: null,
          movesTaken: 0,
          ingameState: 'failed',
        })
        gameStore.setHighlightMovement(undefined)
        break

      default:
        throw new Error(`Unhandled packet type: ${packet.type}`)
    }
  }

  ws.onclose = (event) => {
    console.warn(`WebSocket closed: ${event.reason}`)
    connection.value = 'disconnected'
    reconnectTimeout = setTimeout(connectWebSocket, RECONNECT_INTERVAL)
  }

  ws.onerror = () => {
    ws.close()
  }
}

onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  clearTimeout(reconnectTimeout)
  socket.value?.close()
})
</script>
