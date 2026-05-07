import { Store } from '@tanstack/store'

type Client = {
  ID: string
}

type ClientStore = {
  me: Client | null
  clients: Array<Client>
}

const clientsStore = new Store<ClientStore>({
  me: null,
  clients: [],
})

export type { Client, ClientStore }
export { clientsStore }
