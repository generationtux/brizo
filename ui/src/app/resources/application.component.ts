import { ActivatedRoute }     from '@angular/router'
import { Component }          from '@angular/core'

import { ApplicationService } from './application.service'
import { Pod }                from './pod.component'

@Component({
  selector:    'application',
  templateUrl: './application.component.html',
  providers:   [ApplicationService],
})

export class ApplicationComponent {
  private application: Application

  constructor(private applicationService: ApplicationService, private route: ActivatedRoute) {
    let name: string;
    route.params.subscribe(
      data => name = data['name']
    )
    applicationService.getApplication(name).subscribe(
      data => this.application = data,
      err => console.error('There was an error: ' + err),
      () => console.log('application loaded'),
    )
  }
}

export class Application {
    constructor(
        public id: number,
        public name: string,
        public pods: Array<Pod>
    ){}
}
