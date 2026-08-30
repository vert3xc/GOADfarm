## Quick start

1. **Copy and edit environment variables:**

   Copy `.env_example` to `.env` and edit the variables.

   ```sh
   cp .env_example .env
   ```

2. **Run it with docker compose:**

   ```sh
   docker compose up --build
   ```

3. **Run exploits:**

   ```sh
   python start_sploit.py your_sploit.* --token <team's api token(from docker compose logs)> -u <farm server url>
   ```

   The API key for the farm can be found in docker's logs.

4. **Feeders:**
	"Feeding" is an attack/defense technique that works by exploiting point calculating coefficients used by most A/D systems.
	 Feeding involves putting highly valuable flags to easy-accessible places in order to cut a team's score.
	 To add a feeder write a python script with a flag as its argument that would put the flag to an easy-accessible spot, go to /feed and add a triplet there.
	 The farm's inner mechanisms will automatically run the script depending on the frequency and target you specified.
