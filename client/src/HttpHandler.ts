import axios from "axios"
import { config } from "./config"

/**
 * player can conn as either player 0 or 1.  
 * when invoked, returns a jwt authToken
 */
export async function authenticate(player: 0 | 1){

    const url = config.urls.auth

    try {
        const res = await axios.post(url + "/snl-authenticate", {player})
        const {token} = res.data
        return token
    } catch (error) {
        console.error(error)
    }

}