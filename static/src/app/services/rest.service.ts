import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import {Feature} from '../components/control/control.component';

@Injectable({
  providedIn: 'root'
})
export class RestService {

  constructor(private http: HttpClient) { }

  private featuresURL = environment.host + environment.postfix + '/features';

  public getFeatures(): Observable<any> {
    return this.http.get(`${this.featuresURL}`);
  }

  public createFeature(feature: Feature): Observable<any> {
    return this.http.post(`${this.featuresURL}`, feature);
  }

  public updateFeature(techName: string, displayName: string, description: string, date: Date, active: boolean): Observable<any> {
    return this.http.put(`${this.featuresURL}/${techName}`, {
      displayName,
      expiresOn:   date,
      active:      active.valueOf().toString(),
      description,
    });
  }

  public deleteFeature(name: string): Observable<any> {
    return this.http.delete(`${this.featuresURL}/${name}`);
  }

}
