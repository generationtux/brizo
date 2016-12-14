import { Pod } from './pod.component'

export class Application {
    constructor(
        public id: number,
        public name: string,
        public pods: Array<Pod>
      ){}
}
