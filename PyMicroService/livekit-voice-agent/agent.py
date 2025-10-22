from dotenv import load_dotenv
import aiohttp
import asyncio
import os
import json
from datetime import datetime

from livekit import agents
from livekit.agents import AgentSession, Agent, RoomInputOptions
from livekit.plugins import noise_cancellation, silero
from livekit.plugins.turn_detector.multilingual import MultilingualModel
from livekit.protocol.room import ListParticipantsRequest
from livekit.api import LiveKitAPI

load_dotenv(".env.local")


class Assistant(Agent):
    def __init__(self, user_participant=None, initial_faqs=None) -> None:
        super().__init__(instructions=self.build_instructions(initial_faqs or []))
        self.user_participant = user_participant
        self.last_question = None
        self.last_reply = None
        self.faqs = initial_faqs or []

    @classmethod
    async def create(cls, user_participant):
        faqs = await cls.load_faqs()
        assistant = cls(user_participant=user_participant, initial_faqs=faqs)
        asyncio.create_task(assistant.update_faqs_loop())
        return assistant

    def get_instructions(self) -> str:
        """Override to dynamically build instructions on every turn."""
        return self.build_instructions(self.faqs)

    def build_instructions(self, faqs: list) -> str:
        base = """You are a helpful voice AI assistant.
            You eagerly assist users with their questions by providing information from your extensive knowledge.
            Your responses are concise, to the point, and without any complex formatting or punctuation including emojis, asterisks, or other symbols.
            You are curious, friendly, and have a sense of humor.
            You work in a salon named 'fake salon', you attend the call on behalf of the salon staff.
            Salon stays open from 9 AM to 6 PM, Monday to Saturday.
            On Sundays and public holidays, salon stays open from 8 AM to 12 PM.
            Salon offers services like haircuts, styling, coloring, and treatments.
            Start with a friendly greeting and introduction about salon do not forget to mention salon name.
            *Do respond only with the information provided regarding the salon; do not assume or invent answers.*
            If there is anything you don't know, politely say: "I'm sorry, I don't have that information at the moment. Let me check with my supervisor and get back to you."
            Please don't say anything else apart from the sentence above.
            Below are some recent Q&A pairs for your reference:"""

        if not faqs:
            return base

        faqs_text = "\n\n".join([f"Q: {item['question']}\nA: {item['answer']}" for item in faqs])
        return f"{base}\n\n{faqs_text}"

    async def update_faqs_loop(self):
        """Keep polling FAQs every 5 seconds and update in memory."""
        while True:
            try:
                new_faqs = await self.load_faqs()
                if new_faqs != self.faqs:
                    self.faqs = new_faqs
                    print(f"✅ FAQs updated. Total: {len(self.faqs)}")
            except Exception as e:
                print(f"⚠️ Error updating FAQs: {e}")
            await asyncio.sleep(5)

    @staticmethod
    async def load_faqs():
        async with aiohttp.ClientSession() as session:
            async with session.get("http://localhost:8080/api/v1/queries/faqs") as resp:
                data = await resp.json()
                return data.get("faqs", [])

    def llm_node(self, chat_ctx, tools, model_settings):
        async def custom_stream():
            chat_ctx.instructions = self.get_instructions()
            full_response = ""

            async for chunk in Agent.default.llm_node(self, chat_ctx, tools, model_settings):
                if isinstance(chunk, str):
                    full_response += chunk
                elif hasattr(chunk, "delta") and chunk.delta and chunk.delta.content:
                    full_response += chunk.delta.content
                yield chunk

            print("AI final response:", full_response)

            uncertain_phrases = [
                "i'm not sure",
                "i don't know",
                "I'm sorry, I don't have that information at the moment. Let me check with my supervisor and get back to you",
                "I'm sorry I don't have that information at the moment Let me check with my supervisor and get back to you"
            ]
            print('😭uncertainity')
            if any(p in full_response for p in uncertain_phrases):
                print("Detected uncertain response, triggering fallback to Go backend.")
                await self.fallback_to_go_backend(self.last_question)

            self.last_reply = full_response

        return custom_stream()

    async def fallback_to_go_backend(self, question):
        async with aiohttp.ClientSession() as session:
            payload = {
                "customer_id": self.user_participant.sid,
                "query_text": question[0],
                "answer": "",
                "query_status": 0  # PENDING
            }
            print("🔄 Falling back to Go backend:", payload)

            with open("agent_debug.log", "a", encoding="utf-8") as f:
                timestamp = datetime.now().strftime("[%Y-%m-%d %H:%M:%S]")
                f.write(f"{timestamp} 🔄 Falling back to Go backend: {json.dumps(payload)}\n")

            await session.post("http://localhost:8080/api/v1/queries", json=payload)

    async def on_user_turn_completed(self, turn_ctx, new_message) -> None:
        self.last_question = new_message.content
        print("User asked:", self.last_question)


async def entrypoint(ctx: agents.JobContext):
    room = ctx.room
    lkapi = LiveKitAPI(
        url=os.getenv("LIVEKIT_URL"),
        api_key=os.getenv("LIVEKIT_API_KEY"),
        api_secret=os.getenv("LIVEKIT_API_SECRET")
    )

    participants_resp = await lkapi.room.list_participants(ListParticipantsRequest(room=room.name))
    user_participant = next((p for p in participants_resp.participants if "agent" not in p.identity.lower()), None)

    if user_participant:
        print("👤 Participant Identity:", user_participant.identity)
        print("👤 Participant SID:", user_participant.sid)

    session = AgentSession(
        stt="assemblyai/universal-streaming:en",
        llm="openai/gpt-4.1-mini",
        tts="cartesia/sonic-2:9626c31c-bec5-4cca-baa8-f8ba9e84c8bc",
        vad=silero.VAD.load(),
        turn_detection=MultilingualModel(),
    )

    assistant = await Assistant.create(user_participant=user_participant)

    await session.start(
        room=ctx.room,
        agent=assistant,
        room_input_options=RoomInputOptions(
            noise_cancellation=noise_cancellation.BVC(),
        ),
    )

    await session.generate_reply(
        instructions = assistant.get_instructions()
    )


if __name__ == "__main__":
    agents.cli.run_app(agents.WorkerOptions(entrypoint_fnc=entrypoint))