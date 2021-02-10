import React from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useParams } from 'react-router-dom';
import { useMutation } from '@apollo/client';
import { Button } from '@trussworks/react-uswds';
import { Field, Form, Formik, FormikProps } from 'formik';
import { DateTime } from 'luxon';
import CreateTestDateQuery from 'queries/CreateTestDateQuery';

import Footer from 'components/Footer';
import Header from 'components/Header';
import MainContent from 'components/MainContent';
import PageWrapper from 'components/PageWrapper';
import { ErrorAlert, ErrorAlertMessage } from 'components/shared/ErrorAlert';
import FieldErrorMsg from 'components/shared/FieldErrorMsg';
import FieldGroup from 'components/shared/FieldGroup';
import HelpText from 'components/shared/HelpText';
import Label from 'components/shared/Label';
import { RadioField } from 'components/shared/RadioField';
import { NavLink, SecondaryNav } from 'components/shared/SecondaryNav';
import TextField from 'components/shared/TextField';
import { TestDateForm } from 'types/accessibilityRequest';
import flattenErrors from 'utils/flattenErrors';
import { TestDateValidationSchema } from 'validations/testDateSchema';

import './index.scss';

const TestDate = () => {
  const { t } = useTranslation('accessibility');
  const { accessibilityRequestId } = useParams<{
    accessibilityRequestId: string;
  }>();
  // eslint-disable-next-line no-unused-vars
  const [mutate, mutationResult] = useMutation(CreateTestDateQuery, {
    errorPolicy: 'all'
  });
  const initialValues: TestDateForm = {
    testType: null,
    dateMonth: '',
    dateDay: '',
    dateYear: '',
    score: {
      isPresent: false,
      value: ''
    }
  };

  const onSubmit = (values: TestDateForm) => {
    const testDate = DateTime.fromObject({
      day: Number(values.dateDay),
      month: Number(values.dateMonth),
      year: Number(values.dateYear)
    });

    mutate({
      variables: {
        input: {
          date: testDate,
          score: values.score.isPresent
            ? Math.round(parseFloat(values.score.value) * 10)
            : null,
          testType: values.testType,
          requestID: accessibilityRequestId
        }
      }
    });
  };

  return (
    <PageWrapper className="add-test-date">
      <Header />
      <MainContent className="margin-bottom-5">
        <SecondaryNav>
          <NavLink to={`/508/requests/${accessibilityRequestId}`}>
            {t('tabs.accessibilityRequests')}
          </NavLink>
        </SecondaryNav>
        <div className="grid-container">
          <h1 className="margin-top-6 margin-bottom-5">
            Add a test date for Medicare Change of Office Initiative 1.3
          </h1>
          <div className="grid-row grid-gap-lg">
            <div className="grid-col-9">
              <Formik
                initialValues={initialValues}
                onSubmit={onSubmit}
                validationSchema={TestDateValidationSchema}
                validateOnBlur={false}
                validateOnChange={false}
                validateOnMount={false}
              >
                {(formikProps: FormikProps<TestDateForm>) => {
                  const { errors, setFieldValue, values } = formikProps;
                  const flatErrors = flattenErrors(errors);
                  return (
                    <>
                      {Object.keys(errors).length > 0 && (
                        <ErrorAlert
                          testId="test-date-errors"
                          classNames="margin-top-3"
                          heading="Please check and fix the following"
                        >
                          {Object.keys(flatErrors).map(key => {
                            return (
                              <ErrorAlertMessage
                                key={`Error.${key}`}
                                errorKey={key}
                                message={flatErrors[key]}
                              />
                            );
                          })}
                        </ErrorAlert>
                      )}
                      {mutationResult.error && (
                        <ErrorAlert heading="System error">
                          <ErrorAlertMessage
                            message={mutationResult.error.message}
                            errorKey="system"
                          />
                        </ErrorAlert>
                      )}
                      <Form>
                        <FieldGroup error={!!flatErrors.testType}>
                          <fieldset className="usa-fieldset">
                            <legend className="usa-label margin-bottom-1">
                              What type of test?
                            </legend>
                            <FieldErrorMsg>{flatErrors.testType}</FieldErrorMsg>

                            <Field
                              as={RadioField}
                              checked={values.testType === 'INITIAL'}
                              id="TestDate-TestTypeInitial"
                              name="testType"
                              label="Initial"
                              value="INITIAL"
                            />
                            <Field
                              as={RadioField}
                              checked={values.testType === 'REMEDIATION'}
                              id="TestDate-TestTypeRemediation"
                              name="testType"
                              label="Remediation"
                              value="REMEDIATION"
                            />
                          </fieldset>
                        </FieldGroup>
                        {/* GRT Date Fields */}
                        <FieldGroup
                          error={
                            !!flatErrors.dateMonth ||
                            !!flatErrors.dateDay ||
                            !!flatErrors.dateYear ||
                            !!flatErrors.validDate
                          }
                        >
                          <fieldset className="usa-fieldset margin-top-4">
                            <legend className="usa-label margin-bottom-1">
                              Test date
                            </legend>
                            <HelpText
                              id="TestDate-DateHelp"
                              className="margin-bottom-2"
                            >
                              For example: 4 28 2020
                            </HelpText>
                            <FieldErrorMsg>
                              {flatErrors.dateMonth}
                            </FieldErrorMsg>
                            <FieldErrorMsg>{flatErrors.dateDay}</FieldErrorMsg>
                            <FieldErrorMsg>{flatErrors.dateYear}</FieldErrorMsg>
                            <FieldErrorMsg>
                              {flatErrors.validDate}
                            </FieldErrorMsg>
                            <div
                              className="usa-memorable-date"
                              style={{ marginTop: '-2rem' }}
                            >
                              <div className="usa-form-group usa-form-group--month">
                                <Label htmlFor="TestDate-DateMonth">
                                  Month
                                </Label>
                                <Field
                                  as={TextField}
                                  error={!!flatErrors.dateMonth}
                                  id="TestDate-DateMonth"
                                  maxLength={2}
                                  name="dateMonth"
                                />
                              </div>
                              <div className="usa-form-group usa-form-group--day">
                                <Label htmlFor="TestDate-DateDay">
                                  {t('general:date.day')}
                                </Label>
                                <Field
                                  as={TextField}
                                  error={!!flatErrors.dateDay}
                                  id="TestDate-DateDay"
                                  maxLength={2}
                                  name="dateDay"
                                />
                              </div>
                              <div className="usa-form-group usa-form-group--year">
                                <Label htmlFor="TestDate-DateYear">
                                  {t('general:date.year')}
                                </Label>
                                <Field
                                  as={TextField}
                                  error={!!flatErrors.dateYear}
                                  id="TestDate-DateYear"
                                  maxLength={4}
                                  name="dateYear"
                                />
                              </div>
                            </div>
                          </fieldset>
                        </FieldGroup>
                        <FieldGroup
                          scrollElement="score.isPresent"
                          error={!!flatErrors['score.isPresent']}
                        >
                          <fieldset className="usa-fieldset margin-top-4">
                            <legend className="usa-label margin-bottom-1">
                              Does this test have a test score?
                            </legend>

                            <FieldErrorMsg>
                              {flatErrors['score.isPresent']}
                            </FieldErrorMsg>
                            <Field
                              as={RadioField}
                              checked={values.score.isPresent === false}
                              id="TestDate-HasScoreNo"
                              name="score.isPresent"
                              label="No"
                              onChange={() => {
                                setFieldValue('score.isPresent', false);
                                setFieldValue('score.name', '');
                              }}
                              value={false}
                            />
                            <Field
                              as={RadioField}
                              checked={values.score.isPresent === true}
                              id="TestDate-HasScoreYes"
                              name="score.isPresent"
                              label="Yes"
                              onChange={() => {
                                setFieldValue('score.isPresent', true);
                              }}
                              value
                              aria-describedby="TestDate-ScoreYes"
                            />
                            {values.score.isPresent && (
                              <div className="width-card-lg margin-top-neg-2 margin-left-4 margin-bottom-1">
                                <FieldGroup
                                  scrollElement="score.value"
                                  error={!!flatErrors['score.value']}
                                >
                                  <Label htmlFor="TestDate-ScoreValue">
                                    Test Score
                                  </Label>
                                  <FieldErrorMsg>
                                    {flatErrors['score.value']}
                                  </FieldErrorMsg>
                                  <div className="display-flex">
                                    <div className="width-10">
                                      <Field
                                        as={TextField}
                                        error={!!flatErrors['score.value']}
                                        id="TestDate-ScoreValue"
                                        maxLength={4}
                                        name="score.value"
                                      />
                                    </div>
                                    <div className="bg-black text-white width-5 margin-top-05 display-flex flex-justify-center flex-align-center">
                                      <span className="text-bold">%</span>
                                    </div>
                                  </div>
                                </FieldGroup>
                              </div>
                            )}
                          </fieldset>
                        </FieldGroup>
                        <Button
                          className="margin-top-2 display-block"
                          type="submit"
                        >
                          Add date
                        </Button>
                        <Link
                          to={`/508/requests/${accessibilityRequestId}`}
                          className="margin-top-2 display-block"
                        >
                          Don&apos;t add and go back to the request page
                        </Link>
                      </Form>
                    </>
                  );
                }}
              </Formik>
            </div>
          </div>
        </div>
      </MainContent>
      <Footer />
    </PageWrapper>
  );
};

export default TestDate;
