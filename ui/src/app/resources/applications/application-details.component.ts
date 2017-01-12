import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params }                          from '@angular/router'
import { Component, EventEmitter, OnInit }                 from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/switchMap';

import { ApplicationService } from './application.service'
import { Pod }                from '../pod.component'

@Component({
  selector:    'application',
  templateUrl: './application-details.html',
  providers:   [ApplicationService],
})

export class ApplicationComponent implements OnInit {
  private application: Application;
  public createEnvironmentForm: FormGroup;
  private environments: Environment[];
  public submitted: boolean;
  //private modalActions = new EventEmitter<string>();

  constructor(private applicationService: ApplicationService, private route: ActivatedRoute, private http: Http) {}

  ngOnInit() {
    this.route.params
      .switchMap(
        (params: Params) => this.applicationService.getApplication(params['uuid'])
      ).subscribe(
        data => {
          this.application = data
          this.environments = this.application.environments
        },
        err => console.error('There was an error: ' + err)
      );
    this.createEnvironmentForm = new FormGroup({
      name: new FormControl('', Validators.required)
    });
  }

  save() {
    this.createEnvironment(this.createEnvironmentForm.controls['name'].value)
      .subscribe(
        data => console.log(data),
        err => console.error('There was an error: ' + err),
      );
  }

  createEnvironment(name: string): Observable<Environment> {
    const headers = new Headers({ 'Content-Type': 'application/json' });
    const options = new RequestOptions({ headers: headers });
    const data = { name: name, applicationId: this.application.ID };

    return this.http.post('/api/v1/environments', data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}

export class Application {
  constructor(
    public id:   number,
    public ID:   number,
    public uuid: string,
    public name: string,
    public environments: Environment[],
    public pods: Array<Pod>
  ){}
}

export class Environment {
  constructor(
    public id:             number,
    public uuid:           string,
    public application_id: number,
  ){}
}
